## guides
### install
`helm upgrade --install --create namespace --namespace <namespace> <release> --set config.hostname=<hostname> oci://docker.io/byteplow/idd4 --version <version>`

install from repo: `helm upgrade --install <release> --set config.hostname=<hostname> --set ui.image.tag=latest ./contrib/deployment/chart/`

### add openid connect client
`kubectl exec -it deployments/<release>-hydra -- hydra clients create --endpoint http://localhost:4445 --grant-types authorization_code,refresh_token --response-types code,id_token --scope openid --scope profile --scope email --scope offline --callbacks <https redirect url>`
### register first user
Open `https://<hostname>/self-service/registration/browser?invite=<master invite>` and register.
The master invite is hardcodes as wellknown, but that is subject to change. It will be `configurable` in your values.yaml. 
Or it will be able to be read from a secret. `kubectl get secret <release>-ui -o "jsonpath={.data['masterInvite']}" | base64 --decode`

### backup
The pvc for hydra, kratos and keto must be backed up. And also the secrets for all three. All secret values can also be set in .Values.config. Then only the values.yaml need to be backed up.

The masterInvite key may be backed up. E.g. If some other application is using it.

## diagrams
### deployment
```plantuml
@startuml
node Kratos {
 agent kratospublic
 agent kratosadmin
}

node hydra {
 agent hydrapublic
 agent hydraadmin
}

node keto
node ui

node traefikingress [
traefikingress
---
IngressRoute
]
cloud internet

traefikingress --> hydrapublic : /connect, /oauth2, /userinfo, /.well-known/jwks.json, /.well-known/openid-configuration,
traefikingress --> kratospublic : /sessions, /self-service, /.well-known/ory/webauthn.js, /self-service/registration not for post
traefikingress --> ui : /, /flow, /self-service/registration for post

ui --> kratospublic
ui --> kratosadmin
ui --> hydraadmin
ui --> keto

internet --> traefikingress

file certificate

traefikingress --> certificate : use

@enduml
```

### login flow
```plantuml
@startuml
Browser -> Kratos : get: /self-service/login/browser + (optionaly) redirect_url querry
Browser <- Kratos : redirect: with flow id + set cookie
Browser -> Ui : get: /login?flow=id + cookie
Ui -> Kratos : get login flow with cookie and flow id
Ui <- Kratos : login flow object
Browser <- Ui: rendered login form
note across: User submits data
Browser -> Kratos : post: /self-service/login/browser + form data
alt login successful
Browser <- Kratos : redirect: to redirect_url or default
Browser -> : redirect
else error or additonal login steps like otp
Browser <- Kratos : redirect: with flow id
Browser -> Ui : get: /login?flow=id + cookie
Ui -> Kratos : get login flow with cookie and flow id
Ui <- Kratos : ui login form
Browser <- Ui: rendered login form
note across: User submits data again. Repeat until login success.
end
@enduml
```

### registration flow
```plantuml
@startuml
Browser -> Kratos : get: /self-service/registration/browser + (optionaly) invite querry
Browser <- Kratos : redirect: with flow id + set cookie
Browser -> Ui : get: /registration?flow=id + cookie
Ui -> Kratos : get registration flow with cookie and flow id
Ui <- Kratos : registration flow object
note over Ui: Ui extracts invite from login flow object
alt invite is not empty
Browser <- Ui: rendered login form
else
Browser <- Ui: rendered error message
note across: End here
end
note across: User submits data
Browser -> Ui : post: /self-service/registration/browser + form data
Ui -> Kratos : get registration flow with cookie and flow id
Ui <- Kratos : registration flow object
note over Ui: Ui extracts invite from login flow object
Ui -> keto : check invite
Ui <- keto : invite validity
alt invite is invalid
Browser <- Ui: redirect to ui error page /error
note across: End here
end
Ui -> Kartos : post: /self-service/registration/browser + form data + cookie (proxy forwards http request)
Ui <- Kartos : redirect to /welcome (proxy forwards to browser)
Ui -> keto : invalidate invite
Browser <- Ui : redirect to /welcome (proxy forwards from kartos)
@enduml
```

### settings flow
```plantuml
@startuml
Browser -> Kratos : get: /self-service/settings/browser
Browser <- Kratos : redirect: with flow id + set cookie
Browser -> Ui : get: /settings
Ui -> Kratos : get session with cookie 
Ui <- Kratos : session
Ui -> Kratos : get flow with cookie + flow id
Ui <- Kratos : flow
Browser <- Ui: rendered settings form
note across: User submits data
Browser -> Kratos : post: /self-service/settings/browser + form data
alt valid settings update
Browser -> Kratos : redirect to ui /wellcome
else invalied settings update
Browser -> Kratos : redirect with flow id + set cookie
Browser -> Ui : get: /settings
Ui -> Kratos : get session with cookie 
Ui <- Kratos : session
Ui -> Kratos : get flow with cookie + flow id
Ui <- Kratos : flow
Browser <- Ui: rendered settings form
note across: User submits data again. Repeat untile update is valid
end
@enduml
```