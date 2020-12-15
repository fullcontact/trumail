# Trumail

Trumail is a simple email verification service.

## How it Works

Verifying the deliverability of an email address isn't a very complicated process. In fact, the process Trumail takes to
verify an address is really only half that of sending a standard email transmission and is outlined below...

```
First a TCP connection is formed with the MX server on port 25.

HELO my-domain.com              // We identify ourselves as my-domain.com (set via environment variable)
MAIL FROM: me@my-domain.com     // Set the FROM address being our own
RCPT TO: test-email@example.com // Set the recipient and receive a (200, 500, etc..) from the server
QUIT                            // Cancel the transaction, we have all the info we need
```

As you can see we first form a tcp connection with the mail server on port 25. We then identify ourselves as example.com
and set a reply-to email of admin@example.com (both these are configured via the SOURCE_ADDR environment variable). The
last, and obviously most important step in this process is the RCPT command. This is where, based on the response from
the mail server, we are able to conclude the deliverability of a given email address. A 200 implies a valid inbox and
anything else implies either an error with our connection to the mail server, or a problem with the address requested.

## FC Trumail Service

**Base URL:** http://trumail.profoundis.com

The service is internal. Hence, it can be used only from an internal host or a VPN has to be used.

### Endpoint
There is only 1 endpoint that supports json and XML

#### _GET_ /v1/\<fmt\>/\<email\>
- `fmt` - Format; supports `json` and `xml`
- `email` - The email to be verified

##### Example
```bash
❯ curl http://trumail.profoundis.com/v1/json/chris@fullcontact.com | jq
```
```json
{
  "address": "chris@fullcontact.com",
  "username": "chris",
  "domain": "fullcontact.com",
  "md5Hash": "777b051664a0f5504967ec8f93df2cff",
  "validFormat": true,
  "deliverable": true,
  "fullInbox": false,
  "hostExists": true,
  "catchAll": false,
  "disposable": false,
  "message": "OK",
  "errorDetails": ""
}
```

## Build and Deployment
Current build and deployment process is as follows:
1. Add a tag to the latest commit and push it to this repo.
2. Create a new branch in devops repo and update this tag in [Trumail Ansible configuration](https://github.com/fullcontact/devops/blob/cd11c3711fe84f438e8fa75e29ed1d61d5f39bdb/src/ansible/ansible/roles/drt-trumail/defaults/main.yml#L7)
3. Run the [Ansible Pact Jenkins Job](https://jenkins.eks-cicd.useast1.master.fullcontact.com/job/Ops/job/ansible/job/pact/build) with the below parameters:
  - **sha1**: *The newly created devops repo branch name*
  - **IMAGE_BASE**: *bionic*
  - **image**: *drt-trumail/drt-trumail.json*
4. Once the build is complete, go to the build's console output and copy the ID of the AMI that has been created.
5. Update the AMI ID in the [Spinnaker AMI lookup script in devops repo](https://github.com/fullcontact/devops/blob/cd11c3711fe84f438e8fa75e29ed1d61d5f39bdb/src/spinnaker/scripts/ami_lookup.yaml#L42) and create a PR to merge this branch to the master branch of the devops repo.
6. Once merged, start the manual execution of the [pipeline](https://spinnaker.eks-cicd.useast1.master.fullcontact.com/#/applications/drttrumail/executions) for `drttrumail` application in spinnaker