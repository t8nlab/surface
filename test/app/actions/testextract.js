import { extract } from "@titanpl/surface";

export default function testextract() {
  const steps = [];
  try {
    const targetUrl = "https://mail.google.com";

    // 1. Extract Meta
    steps.push(`Natively extracting Meta/SEO from: ${targetUrl}`);
    const meta = extract.meta(targetUrl);

    // 2. Extract Links
    steps.push(`Natively extracting all unique links from: ${targetUrl}`);
    const links = extract.links(targetUrl);

    return {
      success: true,
      message: "Native extraction verified",
      results: {
        meta,
        linksCount: links.length,
        sampleLinks: links.slice(0, 5)
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
