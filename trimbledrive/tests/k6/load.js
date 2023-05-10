/*
  Load Testing is primarily concerned with assessing the current performance of system in terms of concurrent users or requests per second.

  Load Testing is used to ensure that the application performs satisfactorily when many users access it at the same time.
  Load tests are run regularly and must be run before approving the version for STAGE/PROD release, but not on every PR because they are taking significant time.
  Load tests must be run against a system that are provisioned on pair with PROD from infrastructure point of view.

  Run a load test to:
  - Assess the current performance of system under typical and peak load.
  - Make sure system continue to meet the performance standards as you make changes to system (code and infrastructure).

  Implementation approach:
  Several scenarios are run in parallel that are making random requests to the system with various query parameters and for different data set sizes.
  Each scenario can be run in several variations by controlling following parameters:
  - response compression algo: br, gzip, identity
  - models size: small, medium, large

  The main scenario calls all endpoints with various variations of parameters:
  - all API endpoints
  - full data and with filtering/searching
  - successful and error requests (error requests are extracted to separate scenario)
  - various query parameters that affect the amount of data returned
  - first and last pages in the collection responses
  - small and large pages
*/

import http from 'k6/http';
import { Trend } from 'k6/metrics';
import { intUrl, shared_token } from './constants.js';

//
// Input parameters
//
const token = __ENV.TC_TOKEN ? `${__ENV.TC_TOKEN}` : shared_token;
const baseUrl = __ENV.BASE_URL ? `${__ENV.BASE_URL}` : intUrl;
console.log(baseUrl);

//
// Metrics and configuration
//

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    discardResponseBodies: false,
    scenarios: {
        default: {
            executor: 'ramping-vus',
            stages: [
                { duration: '5m', target: 6 },    // simulate ramp-up of traffic from 1 to target users over 5 minutes
                { duration: '10m', target: 6 },   // stay at target users for 10 minutes
                { duration: '5m', target: 0 },     // ramp-down to 0 users
            ],
        },
        bad_requests: {
            executor: 'ramping-vus',
            exec: 'bad_requests',
            stages: [
                { duration: '5m', target: 3 },    // simulate ramp-up of traffic from 1 to target users over 5 minutes
                { duration: '10m', target: 3 },   // stay at target users for 10 minutes
                { duration: '5m', target: 0 },      // ramp-down to 0 users
            ],
        },
    },

    // thresholds are used to make the per endpoint stats visible in a summary
    thresholds: {
        "http_req_failed": ['rate<0.01'], //  <1% errors
        // rest is added dynamically below
    }
};

options.thresholds[`http_req_duration{name:LandingPage}`] = ['med<50'];
options.thresholds[`http_req_duration{name:Ping}`] = ['med<50'];
options.thresholds[`http_req_duration{name:Health}`] = ['med<500'];
options.thresholds[`http_req_duration{name:NotExists}`] = ['med<500'];
options.thresholds[`http_req_duration{name:SpaceCreate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceCreate.Options}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceCreate.Unauthorized}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceGet.NotFound}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceUpdate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceDelete}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderCreate}`] = ['med<120'];
options.thresholds[`http_req_duration{name:FolderGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderGetVersionNotFound}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderGetVersion}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderGetVersionByMajor}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderUpdate}`] = ['med<120'];
options.thresholds[`http_req_duration{name:FolderDelete}`] = ['med<120'];
options.thresholds[`http_req_duration{name:FolderListChildren}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderListVersions}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FileCreateUpload}`] = ['med<100'];
options.thresholds[`http_req_duration{name:UploadSrcContent}`] = ['med<120'];
options.thresholds[`http_req_duration{name:FileGetUpload}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FileGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FileUpdate}`] = ['med<120'];
options.thresholds[`http_req_duration{name:FileGetByMajorVersion}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FileGetByFullVersion}`] = ['med<100'];
options.thresholds[`http_req_duration{name:ListFileVersions}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FileDelete}`] = ['med<120'];



const folderVersionWaitTime = new Trend('FolderVersionWaitTime', true);
const fileUploadCompletionTime = new Trend('FileUploadCompletionTime', true)


const only200Callback = http.expectedStatuses(200);
const only200or204Callback = http.expectedStatuses(200, 204);
const only201Callback = http.expectedStatuses(201);
const only202Callback = http.expectedStatuses(202);
const only204Callback = http.expectedStatuses(204);
const only400Callback = http.expectedStatuses(400);
const only401Callback = http.expectedStatuses(401);
const only404Callback = http.expectedStatuses(404);
const only405Callback = http.expectedStatuses(405);
const only404Or200Callback = http.expectedStatuses(200, 404);

export default function () {
    const headers = {
        'Authorization': `Bearer ${token}`,
        //'Accept-Encoding': 'gzip', // br, gzip, identity
        'Content-Type': `application/json`,
    };
    const srcContentUploadHeaders = {
        "x-ms-blob-type": "BlockBlob",
        "Content-Type": "text/plain"
    }
    const fileContent = "Test Upload";

    // first request creates an http connection that will be reused across next requests in the iteration, metric for this request includes connection overhead
    http.get(`${baseUrl}`, { responseCallback: only404Callback, tags: { name: `LandingPage` } });

    http.get(`${baseUrl}/v1/app/health/ping`, { responseCallback: only200Callback, tags: { name: `Ping` } });
    http.get(`${baseUrl}/v1/app/health`, { responseCallback: only200Callback, tags: { name: `Health` } });

    const space = http.post(`${baseUrl}/v1/spaces`, '{}', { headers: headers, responseCallback: only201Callback, tags: { name: `SpaceCreate` } });
    const spaceId = space.json().id;
    const rootId = space.json().rootId;

    http.get(`${baseUrl}/v1/spaces/${spaceId}`, { headers: headers, responseCallback: only200Callback, tags: { name: `SpaceGet` } });
    http.patch(`${baseUrl}/v1/spaces/${spaceId}`, '{}', { headers: headers, responseCallback: only200Callback, tags: { name: `SpaceUpdate` } });

    const folder = http.post(`${baseUrl}/v1/spaces/${spaceId}/folders`, `{"name":"folder1","parentFolderId":"${rootId}"}`, { headers: headers, responseCallback: only201Callback, tags: { name: `FolderCreate` } });
    const folderId = folder.json().id;

    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${rootId}/children`, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderListChildren` } });

    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderGet` } });
    http.patch(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, '{}', { headers: headers, responseCallback: only200Callback, tags: { name: `FolderUpdate` } });

    // measure time needed for version availability
    const start = Date.now();
    do {
        const res = http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}/versions/1.0`, { headers: headers, responseCallback: only404Or200Callback, tags: { name: `FolderGetVersionNotFound` } });
        if (res.status === 200) break;
    } while (true);
    folderVersionWaitTime.add(Date.now() - start);
    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}/versions/1.0`, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderGetVersion` } });
    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}/versions/1`, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderGetVersionByMajor` } });

    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}/versions`, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderListVersions` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, null, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderDelete` } });

    const upload = http.post(`${baseUrl}/v1/spaces/${spaceId}/uploads`, `{"name":"file1","parentFolderId":"${rootId}"}`, { headers: headers, responseCallback: only202Callback, tags: { name: `FileCreateUpload` } });
    const uploadId = upload.json().uploadId;
    const srcUploadUrl = upload.json().input.contents[0].url;
    const fileId =  upload.json().input.fileId

    http.put(`${srcUploadUrl}`, `${fileContent}`, { headers: srcContentUploadHeaders, responseCallback: only201Callback, tags: { name: `UploadSrcContent` } })
    // measure time needed for file Upload
    const startTime = Date.now();

    do {
        const resp = http.get(`${baseUrl}/v1/spaces/${spaceId}/uploads/${uploadId}`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileGetUpload` } });
        if (resp.status === 200 && resp.json().status==="DONE") {
            break;
        }
    } while (true);
    fileUploadCompletionTime.add(Date.now() - startTime)

    http.get(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileGet` } });

    http.patch(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}`, `{"name":"updatedFile"}`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileUpdate` } });


    http.get(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}/versions/1`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileGetByMajorVersion` } });

    http.get(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}/versions/1.0`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileGetByFullVersion` } });

    http.get(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}/versions`, { headers: headers, responseCallback: only200Callback, tags: { name: `ListFileVersions` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}/files/${fileId}`, null, { headers: headers, responseCallback: only200Callback, tags: { name: `FileDelete` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, null, { headers: headers, responseCallback: only200Callback, tags: { name: `FolderDelete` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}`, null, { headers: headers, responseCallback: only204Callback, tags: { name: `SpaceDelete` } });
};

export function bad_requests() {
    const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': `application/json`,
    };
    const headersForOptions = {
        'Access-Control-Request-Method': `POST`,
        'Origin': `localhost`,
    };

    // first request creates an http connection that will be reused across next requests in the iteration, metric for this request includes connection overhead
    http.get(`${baseUrl}`, { responseCallback: only404Callback, tags: { name: `LandingPage` } });

    http.options(`${baseUrl}/v1/spaces`, null, { headers: headersForOptions, responseCallback: only200or204Callback, tags: { name: `SpaceCreate.Options` } })

    // client errors
    http.post(`${baseUrl}/v1/spaces`, '{}', { responseCallback: only401Callback, tags: { name: `SpaceCreate.Unauthorized` } });
    http.get(`${baseUrl}/notexists`, { headers: headers, responseCallback: only404Callback, tags: { name: `NotExists` }});
    http.get(`${baseUrl}/v1/spaces/5790d79b-741b-4467-bb49-a7ade0884ae1`, { headers: headers, responseCallback: only404Callback, tags: { name: `SpaceGet.NotFound` } });
};