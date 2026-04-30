import { defineAction, fetch, time } from "@titanpl/native";

export default defineAction((req) => {
  const steps = [];
  try {
    const targetUrl = "https://jsonplaceholder.typicode.com/posts";

    // 1. Test GET with headers and params
    steps.push(`Testing GET with headers and params (Titan Native): ${targetUrl}`);
    const getUrl = `${targetUrl}?_limit=2`;
    const startGet = time.now();
    const getResp = drift(fetch(getUrl, {
      headers: {
        "X-Custom-Header": "Native-Test"
      }
    }));

    const endGet = time.now();

    steps.push(`GET status: ${getResp.status}, OK: ${getResp.ok}`);
    steps.push(`GET data items: ${getResp.length}`);
    steps.push(`GET duration: ${endGet - startGet}ms`);

    // 2. Test POST with JSON data
    steps.push(`Testing POST with JSON data to: ${targetUrl} (Titan Native)`);
    const startPost = time.now();
    const postResp = drift(fetch(targetUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer some-token"
      },
      body: JSON.stringify({
        title: "Native Test",
        body: "Testing native fetch powers",
        userId: 1
      })
    }));
    const endPost = time.now();

    steps.push(`POST status: ${postResp.status}, OK: ${postResp.ok}`);
    steps.push(`POST result ID: ${postResp.id}`);
    steps.push(`POST duration: ${endPost - startPost}ms`);

    return {
      success: true,
      message: "Native fetch verified",
      results: {
        get: {
          status: getResp.status,
          data: getResp,
          duration: endGet - startGet
        },
        post: {
          status: postResp.status,
          data: postResp,
          duration: endPost - startPost
        },
        totalDuration: (endPost - startPost) + (endGet - startGet)
      },
      execution_log: steps
    };
  } catch (err) {
    if (err === "__SUSPEND__" || err.message === "__SUSPEND__") {
      throw err
    }
    return {
      success: false,
      error: err.message,
      execution_log: steps
    }
  }
});