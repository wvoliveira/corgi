import http from 'k6/http';
import { sleep } from 'k6';

const BASE_URL = 'http://localhost:8081'

export const options = {
  stages: [
    { duration: '10s', target: 10 }, // below normal load
    { duration: '40s', target: 10 },
    { duration: '10s', target: 50 }, // spike to 1400 users
    { duration: '40s', target: 50 }, // stay at 1400 for 3 minutes
    { duration: '10s', target: 10 }, // scale down. Recovery stage.
    { duration: '5s', target: 10 },
    { duration: '5s', target: 0 },
  ],
};

export default function () {
  http.get(`${BASE_URL}/api/v1/links`);
  sleep(1);
}
