from datadog import initialize,  api
options = {
        'api_key': '...', ########API Key you can get from Datadog#########
        'app_key': '...'
}

initialize(**options)
api.Monitor.update(8658516,message="@durai {{#is_alert}}@slack-tpaasmonitoring @webhook-api-stg-fs-webhook {{/is_alert}} {{#is_warning}}@webhook-api-stg-fs-webhook {{/is_warning}}{{#is_recovery}}@tpaas_noc@grpmail.trimble.com{{/is_recovery}}")