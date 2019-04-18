package oauthtokenwebhook

import (
	"encoding/json"
	"fmt"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func WebHookHandler(w http.ResponseWriter, r *http.Request, adUsersChan chan string) {
	logrus.Info("Handling oauthtokenwebhook webhook")
	var body []byte
	// Read request body
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	// K8S sends POST request with the admission webhook data,
	// the body can't be empty, but if it is,
	// further processing will be stopped and empty
	// admission response will be sent to K8S API
	if len(body) == 0 {
		sendAdmissionResponse(w)
		logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
		return
	}
	// This object gonna hold actual username
	var oauthToken oauthv1.OAuthAccessToken
	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
		sendAdmissionResponse(w)
		return
	}
	// Try to unmarshal Admission Review raw object to OAuthAccessToken
	if err := json.Unmarshal(ar.Request.Object.Raw, &oauthToken); err != nil {
		logrus.Error("Error during unmarshaling request body")
		logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
		sendAdmissionResponse(w)
		return
	}
	sendAdmissionResponse(w)
	logrus.Infof("Passing user: %s to channel for further processing", oauthToken.UserName)
	// Write AD user to channel for further processing
	adUsersChan <- oauthToken.UserName
	logrus.Info("request is done. . . .")
}

func sendAdmissionResponse(w http.ResponseWriter) {
	/*
	Admission response will always Allowed: true
	Since we don't care about token validation
	The only important thing is an actual trigger for user authentication
	and the OAuthAccessToken object which is includes actual username
	 */
	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse = &v1beta1.AdmissionResponse{Allowed: true, Result: &metav1.Status{Message: "all good"}}
	admissionReview := v1beta1.AdmissionReview{}
	admissionReview.Response = admissionResponse
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		logrus.Errorf("Error during marshaling admissionResponse object: %v", err)
		http.Error(w, fmt.Sprintf("Error during marshaling admissionResponse object: %w", err), http.StatusInternalServerError)
	}
	logrus.Info("Sending response to API server")
	if _, err := w.Write(resp); err != nil {
		logrus.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
