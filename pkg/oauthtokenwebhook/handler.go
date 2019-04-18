package oauthtokenwebhook

import (
	"encoding/json"
	"fmt"
	oauthv1 "github.com/openshift/api/oauth/v1"
	log "github.com/sirupsen/logrus"
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
	log.Info("In handle webhook")
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)
		return
	}
	ar := v1beta1.AdmissionReview{}
	deserializer.Decode(body, nil, &ar);
	var oauthToken oauthv1.OAuthAccessToken
	_ = json.Unmarshal(ar.Request.Object.Raw, &oauthToken)
	log.Info("handling in foo handler 1")
	admissionResponse := v1beta1.AdmissionResponse{Allowed: true, Result: &metav1.Status{Message: "all good"}}
	admissionReview := v1beta1.AdmissionReview{}
	admissionReview.Response = &admissionResponse
	admissionReview.Response.UID = ar.Request.UID
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		log.Error("Error during marshaling admissionResponse object")
		http.Error(w, fmt.Sprintf("Error during marshaling admissionResponse object: %w", err), http.StatusInternalServerError)
	}
	log.Info("Sending response to API server")
	if _, err := w.Write(resp); err != nil {
		log.Error("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}

	//go activedirectory.SyncUserPermissions(oauthToken)
	adUsersChan <- oauthToken.UserName
	log.Info("request is done. . . .")
}
