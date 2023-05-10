/*
  A soak test is used to validate the reliability of the system over a long time.

  Run the soak test to:
  - Verify that system doesn't suffer from bugs or memory leaks, which result in a crash or restart after several hours of operation.
  - Verify that expected application restarts don't lose requests.
  - Find bugs related to race-conditions that appear sporadically.
  - Make sure database doesn't exhaust the allotted storage space and stops.
  - Make sure logs don't exhaust the allotted disk storage.
  - Make sure the external services system depends on don't stop working after a certain amount of requests are executed.

  How to run a soak test:
  - Determine the maximum amount of users the system cam handle.
  - Get the 75-80% of that value
  - Set VUs to that value
  - Run the test in 3 stages. Rump up to the VUs, stay they for 4-12 hours, rump down to 0.
*/


import http from 'k6/http';
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
    stages: [
        { duration: '2m', target: 60 }, // ramp up to X users
        { duration: '3h56m', target: 60 }, // stay at X for ~4 hours
        { duration: '2m', target: 0 }, // scale down. (optional)
    ],

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
options.thresholds[`http_req_duration{name:FolderCreate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderUpdate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderDelete}`] = ['med<100'];

const only200Callback = http.expectedStatuses(200);
const only200or204Callback = http.expectedStatuses(200, 204);
const only201Callback = http.expectedStatuses(201);
const only204Callback = http.expectedStatuses(204);
const only400Callback = http.expectedStatuses(400);
const only401Callback = http.expectedStatuses(401);
const only404Callback = http.expectedStatuses(404);
const only405Callback = http.expectedStatuses(405);

export default function () {
    const headers = {
        'Authorization': `Bearer ${token}`,
        //'Accept-Encoding': 'gzip', // br, gzip, identity
        'Content-Type': `application/json`,
    };

    const headersForOptions = {
        'Access-Control-Request-Method': `POST`,
        'Origin': `localhost`,
    };

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

    http.get(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, { headers: headers,  responseCallback: only200Callback, tags: { name: `FolderGet` } });
    http.patch(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, '{}', { headers: headers,  responseCallback: only200Callback, tags: { name: `FolderUpdate` } });
    http.del(`${baseUrl}/v1/spaces/${spaceId}/folders/${folderId}`, null, { headers: headers,  responseCallback: only200Callback, tags: { name: `FolderDelete` } });

    http.del(`${baseUrl}/v1/spaces/${spaceId}`, null, { headers: headers,  responseCallback: only204Callback, tags: { name: `SpaceDelete` } });

    http.options(`${baseUrl}/v1/spaces`, null, { headers: headersForOptions, responseCallback: only200or204Callback, tags: { name: `SpaceCreate.Options` } })

    // client errors
    http.post(`${baseUrl}/v1/spaces`, '{}', { responseCallback: http.expectedStatuses({ min: 300, max: 401 }), tags: { name: `SpaceCreate.Unauthorized` } });  // TODO: k6 does not detect http status correctly
    http.get(`${baseUrl}/notexists`, { headers: headers, responseCallback: only404Callback, tags: { name: `NotExists` }});
    http.get(`${baseUrl}/v1/spaces/5790d79b-741b-4467-bb49-a7ade0884ae1`, { headers: headers, responseCallback: only404Callback, tags: { name: `SpaceGet.NotFound` } });
};
