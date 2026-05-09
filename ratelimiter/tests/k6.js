import http from "k6/http";

import { Counter } from "k6/metrics";

export const success200 = new Counter("success_200");
export const rejected429 = new Counter("rejected_429");

export const options = {
	vus: 10,
	duration: "100s",
};

export default function () {
	const payload = JSON.stringify({
		entity_id: "123",
		entity_type: "user",
	});

	const params = {
		headers: {
			"Content-Type": "application/json",
		},
	};

	const res = http.post(
		"http://127.0.0.1:8080/api/v1/ratelimit/consume",
		payload,
		params,
	);

	if (res.status === 200) {
		success200.add(1);
	}

	if (res.status === 429) {
		rejected429.add(1);
	}
}
