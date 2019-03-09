package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "regexp"
  "crypto/tls"
  // "os/exec"
  // "strings"
  "time"
  "bytes"
)

// func Run(cmd string, args ...string) string {
//   out, err := exec.Command(cmd, args...).CombinedOutput()
//   if err != nil {
//     fmt.Fprint(os.Stderr, "Failed running: %s %s\n%s\n",
//       cmd, strings.Join(args, " "), err)
//     os.Exit(1)
//   }
//   return string(out)
// }

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func PatchKnativeServiceViaAPI(service string) {
   host := "https://kubernetes.default.svc" // use in cluster
   //host := "https://c3.us-south.containers.cloud.ibm.com:22141" // use local

   dir := "/run/secrets/kubernetes.io/serviceaccount"
   //dir := "/tmp" // use local

   fmt.Println("patching service: " + string(service))
   fmt.Print("kubernetes API server: " + string(host))


   t, err := ioutil.ReadFile(dir + "/token")
   check(err)
   re := regexp.MustCompile(`\r?\n`)
   token := re.ReplaceAllString(string(t), "")
   //fmt.Print(string(token))

   //cert, err := ioutil.ReadFile(dir + "/ca.crt")
   //check(err)
   //fmt.Print(string(cert))

   n, err := ioutil.ReadFile(dir + "/namespace")
   check(err)
   namespace := re.ReplaceAllString(string(n), "")
   //fmt.Print(string(namespace))

   uri := host + "/apis/serving.knative.dev/v1alpha1/namespaces/" + string(namespace) + "/services/" + string(service)
   fmt.Println("URI: " + string(uri))

   timeVal := time.Now()
   timeStr := timeVal.Format("20060102150405")

   jsonBody := []byte(`[{"op":"replace","path":"/spec/runLatest/configuration/build/metadata/annotations/trigger","value":"` + timeStr + `"}]`)
   requestBody := bytes.NewBuffer(jsonBody)
   fmt.Println("Patch: " + string(jsonBody))


   tr := &http.Transport{
     TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
   }
   client := &http.Client{Transport: tr}
   req, err := http.NewRequest("PATCH", uri, requestBody)
   check(err)
   req.Header.Set("Authorization", "Bearer " + token)
   req.Header.Set("Content-Type", "application/json-patch+json")
   res, err := client.Do(req)
   check(err)
   response, err := ioutil.ReadAll(res.Body)
   fmt.Println("Response: " + string(response))

}

func main() {
  ready := true

  service := os.Getenv("SERVICE")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // Wait for our IBM Cloud setup to finish
    for !ready {
      time.Sleep(200 * time.Millisecond)
    }

    msg := map[string]interface{}{}

    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &msg)
    if err != nil {
      fmt.Printf("Error parsing: %s\n\n%s\n", err, string(body))
      return
    }

    body, _ = json.MarshalIndent(msg, "", "  ")
    fmt.Printf("HEADERS:\n%v\nBODY:\n%s\n\n", r.Header, body)

    if msg["action"] != nil {
      fmt.Printf("Got issue event\n")
    } else if msg["hook"] != nil {
      fmt.Printf("Got hook event\n")
    } else if msg["pusher"] != nil {
      fmt.Printf("Got push event\n")
      PatchKnativeServiceViaAPI(service)
    } else {
      fmt.Printf("Unknown event\n")
    }
  })

  fmt.Print("Listening on port 8080\n")
  http.ListenAndServe(":8080", nil)
}
