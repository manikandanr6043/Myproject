/*
  Spike testing is a variation of a stress testing, but it does not gradually increase the load. Instead it spikes to extreme load over a very short period of time.
  While a stress test allows the SUT (System Under Test) to gradually scale up its infrastructure, a spike test does not.

  Run the test to:
  - Determine how system performs under a sudden surge of traffic.
  - Determine if system recovers once the traffic has subsided.

  Success is based on expectations. Systems generally react in 4 different ways:

  - Excellent: system performance is not degraded during the surge of traffic. Response time is similar during low traffic and high traffic.
  - Good: Response time is slower, but the system does not produce any errors. All requests are handled.
  - Poor: System produces errors during the surge of traffic, but recovers to normal after the traffic subsides.
  - Bad: System crashes, and does not recover after the traffic has subsided.
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
        { duration: '10s', target: 30 },  // below normal load
        { duration: '1m', target: 30 },
        { duration: '10s', target: 500 }, // spike to overload
        { duration: '3m', target: 500 },  // stay at overload for 3 minutes
        { duration: '10s', target: 30 },  // scale down. Recovery stage.
        { duration: '3m', target: 30 },
        { duration: '10s', target: 0 },
      ],

    // thresholds are used to make the per endpoint stats visible in a summary
    thresholds: {
        // TODO
    }
};

options.thresholds[`http_req_duration{name:LandingPage}`] = ['med<50'];
options.thresholds[`http_req_duration{name:Ping}`] = ['med<50'];
options.thresholds[`http_req_duration{name:Health}`] = ['med<500'];
options.thresholds[`http_req_duration{name:SpaceCreate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceUpdate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:SpaceDelete}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderCreate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderGet}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderUpdate}`] = ['med<100'];
options.thresholds[`http_req_duration{name:FolderDelete}`] = ['med<100'];

const only200Callback = http.expectedStatuses(200);
const only201Callback = http.expectedStatuses(201);
const only204Callback = http.expectedStatuses(204);
const only404Callback = http.expectedStatuses(404);

export default function() {
    const headers = {
        'Authorization': `Bearer ${token}`,
        //'Accept-Encoding': 'gzip', // br, gzip, identity
        'Content-Type': `application/json`,
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
};
