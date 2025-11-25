import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// собственный SLI по успешности
export const technical_failures = new Rate('technical_failures');

export const options = {
    scenarios: {
        constant_load: {
            executor: 'constant-arrival-rate',
            rate: 5,
            timeUnit: '1s',
            duration: '60s',
            preAllocatedVUs: 10,
            maxVUs: 20,
        },
    },
    thresholds: {
        technical_failures: ['rate<0.001'], // <= 0.1% технических фейлов
        http_req_duration: ['p(95)<300'],   // p95 < 300ms
    },
};


export default function () {
    const userIndex = Math.floor(Math.random() * 200) + 1;
    const userId = `u${userIndex}`;

    const res = http.get(`${BASE_URL}/users/getReview?user_id=${userId}`);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    // считаем техническим фейлом только сетевые/5xx
    technical_failures.add(res.error !== '' || res.status >= 500);

    // 10% запросов делаем reassign
    if (Math.random() < 0.1) {
        const prId = `pr-${Math.floor(Math.random() * 100) + 1}`;
        const oldReviewer = `u${Math.floor(Math.random() * 200) + 1}`;

        const payload = JSON.stringify({
            pull_request_id: prId,
            old_user_id: oldReviewer,
        });

        const res2 = http.post(
            `${BASE_URL}/pullRequest/reassign`,
            payload,
            { headers: { 'Content-Type': 'application/json' } },
        );

        // 200/404/409 — ок
        check(res2, {
            'reassign technical success': (r) =>
                r.status === 200 || r.status === 409 || r.status === 404,
        });

        // 5xx или сетевые ошибки - фейл
        technical_failures.add(res2.error !== '' || res2.status >= 500);
    }

    sleep(0.1);
}
