# Shortening-URL

## issue
* keys table have two types data and unused key number should > used key number
  * concurrent create, select key need to filter used key will slow
* cronjob generate key number should over use number
* key generate service should separate with web server  
* do not have cache