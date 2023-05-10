// Smoke test is a regular load test, configured for minimal load.
// You want to run a smoke test as a sanity check every time you modify the code. Smoke tests that are run as part of CI pipeline as a quality gate.
// Duration of the test is minimal to remove any barriers to run it often.
//
// The main goals:
// * fail build fast if something is obviously wrong.
// * verify that the system doesn't throw any errors when under minimal load.
//
// Implementation approach:
// Test API endpoints in isolation in ideal conditions without concurrency and stress.
// Expected to be run on temporary CI stacks that are not provisioned with high level infrastructure.

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
    vus: 1,
    iterations: 50,
    //duration: '180s',
    thresholds: {
        "http_req_failed": ['rate<0.01'], //  <1% errors
        // rest is added dynamically below
    }
};

options.thresholds[`http_req_duration{name:LandingPage}`] = ['med<200'];
options.thresholds[`http_req_duration{name:Ping}`] = ['med<200'];
options.thresholds[`http_req_duration{name:Health}`] = ['med<1000'];
options.thresholds[`http_req_duration{name:SpaceCreate}`] = ['med<200'];
options.thresholds[`http_req_duration{name:SpaceGet}`] = ['med<200'];
options.thresholds[`http_req_duration{name:SpaceUpdate}`] = ['med<200'];
options.thresholds[`http_req_duration{name:SpaceDelete}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderCreate}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderGet}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderGetVersionNotFound}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderGetVersion}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderGetVersionByMajor}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderUpdate}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderDelete}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderListChildren}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FolderListVersions}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileCreateUpload}`] = ['med<200'];
options.thresholds[`http_req_duration{name:UploadSrcContent}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileGetUpload}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileGet}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileUpdate}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileGetByMajorVersion}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileGetByFullVersion}`] = ['med<200'];
options.thresholds[`http_req_duration{name:ListFileVersions}`] = ['med<200'];
options.thresholds[`http_req_duration{name:FileDelete}`] = ['med<200'];

const folderVersionWaitTime = new Trend('FolderVersionWaitTime', true);
const fileUploadCompletionTime = new Trend('FileUploadCompletionTime', true)
const only200Callback = http.expectedStatuses(200);
const only201Callback = http.expectedStatuses(201);
const only202Callback = http.expectedStatuses(202);
const only204Callback = http.expectedStatuses(204);
const only404Callback = http.expectedStatuses(404);
const only404Or200Callback = http.expectedStatuses(200,404);

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

    http.get(`${baseUrl}/v1/spaces/${spaceId}/uploads/${uploadId}`, { headers: headers, responseCallback: only200Callback, tags: { name: `FileGetUpload` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}`, null, { headers: headers, responseCallback: only204Callback, tags: { name: `SpaceDelete` } });
};
