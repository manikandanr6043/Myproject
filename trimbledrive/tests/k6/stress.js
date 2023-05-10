// Stress testing is a type of load testing used to determine the limits of the system.
// The purpose of this test is to verify stability and reliability of the system under extreme conditions.
//
// Run the test to:
// * Determine how system behaves under extreme conditions.
// * Determine what is the maximum capacity of system in terms of users or throughput
// * Determine the breaking point of system and its failure mode
// * Determine if system recovers without manual intervention after the stress test is over


import http from 'k6/http';
import { Counter, Rate } from 'k6/metrics';
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
export const errorCount = new Counter("errors");
export const errorRate = new Rate("error_rate");

export const options = {
    insecureSkipTLSVerify: true,
    noConnectionReuse: false,
    discardResponseBodies: false,
    stages: [
        { duration: '2m', target: 50 },     // below normal load
        { duration: '5m', target: 50 },
        { duration: '2m', target: 100 },    // normal load
        { duration: '5m', target: 100 },
        { duration: '2m', target: 150 },    // around the breaking point
        { duration: '5m', target: 150 },
        { duration: '2m', target: 200 },    // beyond the breaking point
        { duration: '5m', target: 200 },
        { duration: '10m', target: 0 },     // scale down. Recovery stage.
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

export default function () {
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
