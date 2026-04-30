import { time } from "@titanpl/native";
import { http } from "@titanpl/surface";

export default function testhttp() {
  const steps = [];
  try {
    const targetUrl = "https://jsonplaceholder.typicode.com/posts";

    // 1. Test GET with headers and params
    steps.push(`Testing GET with headers and params: ${targetUrl}`);
    const startGet = time.now();
    const getResp = http.get(targetUrl, {
      params: { _limit: 2 },
      headers: {
        "X-Custom-Header": "Surface-Test"
      }
    });
    const endGet = time.now();
    
    steps.push(`GET status: ${getResp.status}, OK: ${getResp.ok}`);
    steps.push(`GET data items: ${getResp.data.length}`);
    steps.push(`GET duration: ${endGet - startGet}ms`);

    // 2. Test POST with JSON data
    steps.push(`Testing POST with JSON data to: ${targetUrl}`);
    const startPost = time.now();
    const postResp = http.post(targetUrl, {
      title: "Surface Test",
      body: "Testing native axios-like powers",
      userId: 1
    }, {
      headers: {
        "Authorization": "Bearer some-token"
      }
    });
    const endPost = time.now();
    
    steps.push(`POST status: ${postResp.status}, OK: ${postResp.ok}`);
    steps.push(`POST result ID: ${postResp.data.id}`);
    steps.push(`POST duration: ${endPost - startPost}ms`);

    return {
      success: true,
      message: "Axios-like HTTP powers verified natively",
      results: {
        get: {
          status: getResp.status,
          headers: getResp.headers,
          data: getResp.data,
          duration: endGet - startGet
        },
        post: {
          status: postResp.status,
          data: postResp.data,
          duration: endPost - startPost
        },
        totalDuration: (endPost - startPost) + (endGet - startGet)
      },
      execution_log: steps
    };
  } catch (err) {
    return {
      success: false,
      error: err.message,
      execution_log: steps
    }
  }
}
