import http from 'k6/http';
import { sleep } from 'k6';

const BASE_URL = 'http://127.0.0.1:8081'

export const options = {
  stages: [
    { duration: '10s', target: 1 },
    { duration: '20s', target: 2 },
    { duration: '30s', target: 3 }, // spike to 200 users
    { duration: '50s', target: 4 }, // stay at 300 for 2 minutes
    { duration: '30s', target: 3 }, // scale down. Recovery stage.
    { duration: '20s', target: 2 },
    { duration: '10s', target: 1 },
  ],
};

export default function () {
  const res = http.get(`${BASE_URL}/api/health`);
}
