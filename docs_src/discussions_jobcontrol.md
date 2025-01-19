# Job Contron API Discussions

## Create Trigger Example

Reference: [Trying to Create DB Import Job via API but it's not updating the connection](https://forums.saviynt.com/t5/identity-governance/trying-to-create-db-import-job-via-api-but-it-s-not-updating-the/m-p/24465/highlight/true)

```json
{
  "triggers": [
    {
      "valueMap": {
        "CONNECTION": "DBConnect",
        "connectiontype": "DB",
        "connectionid": 7
      },
      "name": "DBAccessImport",
      "group": "GRAILS_JOBS",
      "jobName": "AccountsImportFullJob",
      "jobGroup": "DATABASE",
      "cronExp": "0 15 10 * * ? 2099"
    }
  ]
}
```

## Create Trigger Chain Example

Reference: [How to make API Call for Trigger Chain Job](https://forums.saviynt.com/t5/saviynt-knowledge-base/how-to-make-api-call-for-trigger-chain-job/ta-p/116827)

```json
{
  "jobgroup": "utility",
  "triggername": "<Chain Job Name>",
  "jobname": "TriggerChainJob",
  "createJobIfDoesNotExist": "false",
  "valueMap": {
    "savtriggerorderform": "Jobname1,Jobname2,JobName3",
    "onFailureForm": "Stop"
  }
}
```