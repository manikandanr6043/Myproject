from datadog import initialize,  api
options = {
        'api_key': '....',
        'app_key': '....'
}

initialize(**options)
api.Monitor.update(8658,message="@du {{#is_alert}}@slack-ing @webhook-api-stg- {{/is_alert}} {{#is_warning}}@webhook-api- {{/is_warning}}{{#is_recovery}}@tpaas_noc@grp{{/is_recovery}}")