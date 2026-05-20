import { useEffect, useMemo, useState } from "react";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "";

function App() {
  const [status, setStatus] = useState("checking");
  const [details, setDetails] = useState(null);

  useEffect(() => {
    const controller = new AbortController();

    async function checkApi() {
      try {
        const response = await fetch(`${API_BASE_URL}/healthz`, {
          signal: controller.signal,
        });
        const payload = await response.json();
        setStatus(response.ok ? "online" : "degraded");
        setDetails(payload);
      } catch (error) {
        if (error.name !== "AbortError") {
          setStatus("offline");
          setDetails(null);
        }
      }
    }

    checkApi();

    return () => controller.abort();
  }, []);

  const statusLabel = useMemo(() => {
    if (status === "online") return "API online";
    if (status === "degraded") return "API degraded";
    if (status === "offline") return "API offline";
    return "Checking API";
  }, [status]);

  return (
    <main className="app-shell">
      <section className="hero">
        <p className="eyebrow">StartTech production dashboard</p>
        <h1>Full-stack delivery pipeline</h1>
        <p className="lede">
          React on S3 and CloudFront, Go on EC2 Auto Scaling, Redis on ElastiCache,
          MongoDB Atlas, and deployment automation through GitHub Actions.
        </p>
      </section>

      <section className="status-grid" aria-label="Service status">
        <article className="status-card">
          <span className={`status-dot ${status}`} aria-hidden="true" />
          <div>
            <h2>{statusLabel}</h2>
            <p>{details?.service ?? "starttech-api"}</p>
          </div>
        </article>

        <article className="status-card">
          <strong>{details?.environment ?? "prod"}</strong>
          <p>Environment</p>
        </article>

        <article className="status-card">
          <strong>{details?.uptime_seconds ?? "--"}</strong>
          <p>Uptime seconds</p>
        </article>
      </section>
    </main>
  );
}

export default App;
