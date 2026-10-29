// import React, { useEffect, useState } from "react";
// import { useParams } from "react-router-dom";
// import api from "../api";
// import { Finding, Review } from "../types";
// const ReviewPage = () => {
//   const { reviewId } = useParams();
//   const [review, setReview] = useState<Review | null>(null);
//   const [findings, setFindings] = useState<Finding[]>([]);
//   useEffect(() => {
//     const load = async () => {
//       const rr = await api.get<Review>(`/api/reviews/${reviewId}`);
//       setReview(rr.data);
//       if (rr.data.status === "completed")
//         setFindings(
//           (await api.get<Finding[]>(`/api/reviews/${reviewId}/findings`)).data,
//         );
//     };
//     load();
//     const id = window.setInterval(load, 3000);
//     return () => window.clearInterval(id);
//   }, [reviewId]);
//   if (!review) return <p>Loading review...</p>;
//   return (
//     <>
//       <h1>Review report</h1>
//       <div className="card p-3 mb-4">
//         <div className="row">
//           <div className="col-md-8">
//             <strong>{review.sourceName}</strong>
//             <p className="mb-0">
//               {review.summary || review.errorMessage || "Review is processing."}
//             </p>
//           </div>
//           <div className="col-md-4 text-md-end">
//             <span className={`badge status-${review.status}`}>
//               {review.status}
//             </span>
//           </div>
//         </div>
//       </div>
//       <div className="row g-2 mb-4">
//         {[
//           ["Critical", review.criticalCount, "danger"],
//           ["High", review.highCount, "warning"],
//           ["Medium", review.mediumCount, "info"],
//           ["Low", review.lowCount, "secondary"],
//         ].map(([name, count, color]) => (
//           <div className="col-6 col-md-3" key={String(name)}>
//             <div className={`card border-${color} p-3`}>
//               <small>{name}</small>
//               <strong className="fs-3">{count}</strong>
//             </div>
//           </div>
//         ))}
//       </div>
//       {findings.map((f) => (
//         <div className="card mb-3" key={f.id}>
//           <div className="card-header d-flex justify-content-between">
//             <strong>{f.title}</strong>
//             <span className={`severity severity-${f.severity.toLowerCase()}`}>
//               {f.severity}
//             </span>
//           </div>
//           <div className="card-body">
//             <p>
//               <strong>
//                 {f.filePath}:{f.lineNumber}
//               </strong>{" "}
//               · {f.category} · {f.ruleId}
//             </p>
//             <p>{f.description}</p>
//             <div className="recommendation">
//               <strong>Recommended fix:</strong> {f.recommendation}
//             </div>
//             {f.codeExcerpt && (
//               <pre className="mt-3">
//                 <code>{f.codeExcerpt}</code>
//               </pre>
//             )}
//           </div>
//         </div>
//       ))}
//       {review.status === "completed" && findings.length === 0 && (
//         <div className="alert alert-success">
//           No built-in rule findings were detected.
//         </div>
//       )}
//     </>
//   );
// };
// export default ReviewPage;

import React, { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import api from "../api";
import { Finding, Review } from "../types";
import "./ReviewPage.css";

type SeverityType = "critical" | "high" | "medium" | "low";

interface SeverityCardProps {
    name: string;
    count: number;
    severity: SeverityType;
    icon: string;
}

const SeverityCard = ({
    name,
    count,
    severity,
    icon,
}: SeverityCardProps) => {
    return (
        <div className={`severity-summary-card ${severity}`}>
            <div className="severity-summary-top">
                <div className={`severity-summary-icon ${severity}`}>
                    <i className={`bi ${icon}`}></i>
                </div>

                <span className={`severity-indicator ${severity}`}>
                    {name}
                </span>
            </div>

            <div className="severity-summary-count">{count}</div>

            <div className="severity-summary-label">
                {count === 1 ? "Finding detected" : "Findings detected"}
            </div>
        </div>
    );
};

const ReviewPage = () => {
    const { reviewId } = useParams();

    const [review, setReview] = useState<Review | null>(null);
    const [findings, setFindings] = useState<Finding[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [refreshing, setRefreshing] = useState(false);

    const loadReview = useCallback(
        async (showRefreshLoader = false) => {
            if (!reviewId) {
                setError("Review ID was not found.");
                setLoading(false);
                return;
            }

            try {
                if (showRefreshLoader) {
                    setRefreshing(true);
                }

                const reviewResponse = await api.get<Review>(
                    `/api/reviews/${reviewId}`
                );

                const reviewData = reviewResponse.data;

                setReview(reviewData);
                setError("");

                if (reviewData.status.toLowerCase() === "completed") {
                    const findingsResponse = await api.get<Finding[]>(
                        `/api/reviews/${reviewId}/findings`
                    );

                    setFindings(findingsResponse.data);
                }
            } catch (requestError) {
                console.error("Unable to load review:", requestError);

                setError(
                    "Unable to load the review report. Please try again."
                );
            } finally {
                setLoading(false);
                setRefreshing(false);
            }
        },
        [reviewId]
    );

    useEffect(() => {
        loadReview();

        const intervalId = window.setInterval(() => {
            loadReview();
        }, 3000);

        return () => {
            window.clearInterval(intervalId);
        };
    }, [loadReview]);

    const getStatusClass = (status: string) => {
        const normalizedStatus = status.toLowerCase();

        if (normalizedStatus === "completed") {
            return "review-report-status completed";
        }

        if (normalizedStatus === "failed") {
            return "review-report-status failed";
        }

        if (
            normalizedStatus === "running" ||
            normalizedStatus === "processing"
        ) {
            return "review-report-status running";
        }

        if (
            normalizedStatus === "pending" ||
            normalizedStatus === "queued"
        ) {
            return "review-report-status pending";
        }

        return "review-report-status default";
    };

    const getStatusIcon = (status: string) => {
        const normalizedStatus = status.toLowerCase();

        if (normalizedStatus === "completed") {
            return "bi-check-circle-fill";
        }

        if (normalizedStatus === "failed") {
            return "bi-x-circle-fill";
        }

        if (
            normalizedStatus === "running" ||
            normalizedStatus === "processing"
        ) {
            return "bi-arrow-repeat";
        }

        return "bi-clock-fill";
    };

    const getSeverityIcon = (severity: string) => {
        const normalizedSeverity = severity.toLowerCase();

        if (normalizedSeverity === "critical") {
            return "bi-exclamation-octagon-fill";
        }

        if (normalizedSeverity === "high") {
            return "bi-exclamation-triangle-fill";
        }

        if (normalizedSeverity === "medium") {
            return "bi-exclamation-circle-fill";
        }

        return "bi-info-circle-fill";
    };

    if (loading) {
        return (
            <div className="review-report-page">
                <div className="review-report-loading">
                    <div className="review-loading-spinner"></div>

                    <h3>Loading review report</h3>

                    <p>
                        Retrieving code analysis results and findings.
                    </p>
                </div>
            </div>
        );
    }

    if (error && !review) {
        return (
            <div className="review-report-page">
                <div className="review-report-container">
                    <div className="review-error-state">
                        <div className="review-error-icon">
                            <i className="bi bi-exclamation-triangle"></i>
                        </div>

                        <h2>Unable to load review</h2>

                        <p>{error}</p>

                        <button
                            type="button"
                            className="review-primary-button"
                            onClick={() => loadReview(true)}
                        >
                            <i className="bi bi-arrow-clockwise"></i>
                            Try again
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    if (!review) {
        return null;
    }

    const normalizedStatus = review.status.toLowerCase();
    const isCompleted = normalizedStatus === "completed";
    const isFailed = normalizedStatus === "failed";

    return (
        <div className="review-report-page">
            <div className="review-report-container">
                <div className="review-report-navigation">
                    <Link
                        to="/dashboard"
                        className="review-back-link"
                    >
                        <i className="bi bi-arrow-left"></i>
                        Back to projects
                    </Link>
                </div>

                <header className="review-report-header">
                    <div>
                        <div className="review-report-label">
                            <i className="bi bi-file-earmark-code"></i>
                            CODE REVIEW REPORT
                        </div>

                        <h1>Review report</h1>

                        <p>
                            Review detected issues, severity levels, affected
                            source files, and recommended fixes.
                        </p>
                    </div>

                    <button
                        type="button"
                        className="review-refresh-button"
                        onClick={() => loadReview(true)}
                        disabled={refreshing}
                    >
                        <i
                            className={`bi bi-arrow-clockwise ${
                                refreshing ? "rotating" : ""
                            }`}
                        ></i>

                        {refreshing ? "Refreshing..." : "Refresh"}
                    </button>
                </header>

                {error && (
                    <div className="review-report-message error">
                        <i className="bi bi-exclamation-circle-fill"></i>
                        <span>{error}</span>
                    </div>
                )}

                <section className="review-overview-card">
                    <div className="review-overview-icon">
                        <i className="bi bi-code-square"></i>
                    </div>

                    <div className="review-overview-content">
                        <div className="review-overview-heading">
                            <div>
                                <span className="review-overview-caption">
                                    REVIEW SOURCE
                                </span>

                                <h2>{review.sourceName}</h2>
                            </div>

                            <span className={getStatusClass(review.status)}>
                                <i
                                    className={`bi ${getStatusIcon(
                                        review.status
                                    )}`}
                                ></i>

                                {review.status}
                            </span>
                        </div>

                        <p className="review-overview-summary">
                            {review.summary ||
                                review.errorMessage ||
                                "Review is currently processing. Results will appear automatically when analysis completes."}
                        </p>

                        {!isCompleted && !isFailed && (
                            <div className="review-processing">
                                <div className="review-processing-animation">
                                    <span></span>
                                    <span></span>
                                    <span></span>
                                </div>

                                <div>
                                    <strong>
                                        Analysis in progress
                                    </strong>

                                    <p>
                                        This page automatically refreshes while
                                        the review is running.
                                    </p>
                                </div>
                            </div>
                        )}
                    </div>
                </section>

                <section className="review-summary-section">
                    <div className="review-section-heading">
                        <div>
                            <span className="review-section-label">
                                ANALYSIS SUMMARY
                            </span>

                            <h2>Findings by severity</h2>
                        </div>

                        {isCompleted && (
                            <div className="review-summary-completed">
                                <i className="bi bi-check-circle-fill"></i>
                                Analysis complete
                            </div>
                        )}
                    </div>

                    <div className="severity-summary-grid">
                        <SeverityCard
                            name="Critical"
                            count={review.criticalCount}
                            severity="critical"
                            icon="bi-exclamation-octagon-fill"
                        />

                        <SeverityCard
                            name="High"
                            count={review.highCount}
                            severity="high"
                            icon="bi-exclamation-triangle-fill"
                        />

                        <SeverityCard
                            name="Medium"
                            count={review.mediumCount}
                            severity="medium"
                            icon="bi-exclamation-circle-fill"
                        />

                        <SeverityCard
                            name="Low"
                            count={review.lowCount}
                            severity="low"
                            icon="bi-info-circle-fill"
                        />
                    </div>
                </section>

                <section className="review-findings-section">
                    <div className="review-section-heading findings-heading">
                        <div>
                            <span className="review-section-label">
                                DETAILED FINDINGS
                            </span>

                            <h2>Code review findings</h2>

                            <p>
                                Review each detected issue and the recommended
                                remediation.
                            </p>
                        </div>

                        {isCompleted && findings.length > 0 && (
                            <div className="review-findings-count">
                                {findings.length}
                                <span>
                                    {findings.length === 1
                                        ? "Finding"
                                        : "Findings"}
                                </span>
                            </div>
                        )}
                    </div>

                    {findings.length > 0 && (
                        <div className="review-findings-list">
                            {findings.map((finding, index) => {
                                const severity =
                                    finding.severity.toLowerCase();

                                return (
                                    <article
                                        className={`review-finding-card severity-${severity}`}
                                        key={finding.id}
                                    >
                                        <div className="finding-card-header">
                                            <div className="finding-heading-content">
                                                <div
                                                    className={`finding-severity-icon ${severity}`}
                                                >
                                                    <i
                                                        className={`bi ${getSeverityIcon(
                                                            finding.severity
                                                        )}`}
                                                    ></i>
                                                </div>

                                                <div>
                                                    <span className="finding-number">
                                                        Finding {index + 1}
                                                    </span>

                                                    <h3>{finding.title}</h3>
                                                </div>
                                            </div>

                                            <span
                                                className={`finding-severity-badge ${severity}`}
                                            >
                                                {finding.severity}
                                            </span>
                                        </div>

                                        <div className="finding-card-body">
                                            <div className="finding-metadata">
                                                <div className="finding-meta-item">
                                                    <span className="finding-meta-icon">
                                                        <i className="bi bi-file-earmark-code"></i>
                                                    </span>

                                                    <div>
                                                        <span className="finding-meta-label">
                                                            File
                                                        </span>

                                                        <strong>
                                                            {finding.filePath}
                                                        </strong>
                                                    </div>
                                                </div>

                                                <div className="finding-meta-item">
                                                    <span className="finding-meta-icon">
                                                        <i className="bi bi-list-ol"></i>
                                                    </span>

                                                    <div>
                                                        <span className="finding-meta-label">
                                                            Line
                                                        </span>

                                                        <strong>
                                                            {finding.lineNumber}
                                                        </strong>
                                                    </div>
                                                </div>

                                                <div className="finding-meta-item">
                                                    <span className="finding-meta-icon">
                                                        <i className="bi bi-tag"></i>
                                                    </span>

                                                    <div>
                                                        <span className="finding-meta-label">
                                                            Category
                                                        </span>

                                                        <strong>
                                                            {finding.category}
                                                        </strong>
                                                    </div>
                                                </div>

                                                <div className="finding-meta-item">
                                                    <span className="finding-meta-icon">
                                                        <i className="bi bi-braces"></i>
                                                    </span>

                                                    <div>
                                                        <span className="finding-meta-label">
                                                            Rule
                                                        </span>

                                                        <strong>
                                                            {finding.ruleId}
                                                        </strong>
                                                    </div>
                                                </div>
                                            </div>

                                            <div className="finding-description">
                                                <h4>
                                                    <i className="bi bi-search"></i>
                                                    What was detected
                                                </h4>

                                                <p>{finding.description}</p>
                                            </div>

                                            <div className="finding-recommendation">
                                                <div className="finding-recommendation-icon">
                                                    <i className="bi bi-lightbulb-fill"></i>
                                                </div>

                                                <div>
                                                    <span>
                                                        RECOMMENDED FIX
                                                    </span>

                                                    <p>
                                                        {
                                                            finding.recommendation
                                                        }
                                                    </p>
                                                </div>
                                            </div>

                                            {finding.codeExcerpt && (
                                                <div className="finding-code-section">
                                                    <div className="finding-code-header">
                                                        <div>
                                                            <span className="code-window-dot red"></span>
                                                            <span className="code-window-dot yellow"></span>
                                                            <span className="code-window-dot green"></span>
                                                        </div>

                                                        <span>
                                                            <i className="bi bi-code-slash"></i>
                                                            Code excerpt
                                                        </span>

                                                        <span className="code-line-label">
                                                            Line{" "}
                                                            {finding.lineNumber}
                                                        </span>
                                                    </div>

                                                    <pre>
                                                        <code>
                                                            {
                                                                finding.codeExcerpt
                                                            }
                                                        </code>
                                                    </pre>
                                                </div>
                                            )}
                                        </div>
                                    </article>
                                );
                            })}
                        </div>
                    )}

                    {isCompleted && findings.length === 0 && (
                        <div className="review-clean-state">
                            <div className="review-clean-icon">
                                <i className="bi bi-shield-check"></i>
                            </div>

                            <div>
                                <span className="review-clean-label">
                                    REVIEW COMPLETE
                                </span>

                                <h3>No findings detected</h3>

                                <p>
                                    No built-in rule findings were detected
                                    during this review.
                                </p>
                            </div>
                        </div>
                    )}

                    {!isCompleted && !isFailed && (
                        <div className="review-waiting-state">
                            <div className="review-waiting-icon">
                                <i className="bi bi-cpu"></i>
                            </div>

                            <h3>Analysis is still running</h3>

                            <p>
                                Detailed findings will appear here when the code
                                review completes.
                            </p>

                            <div className="review-analysis-progress">
                                <span></span>
                            </div>
                        </div>
                    )}

                    {isFailed && (
                        <div className="review-failed-state">
                            <div className="review-failed-icon">
                                <i className="bi bi-x-circle"></i>
                            </div>

                            <h3>Review failed</h3>

                            <p>
                                {review.errorMessage ||
                                    "The code review could not be completed."}
                            </p>

                            <Link
                                to="/dashboard"
                                className="review-primary-link"
                            >
                                <i className="bi bi-arrow-left"></i>
                                Back to projects
                            </Link>
                        </div>
                    )}
                </section>
            </div>
        </div>
    );
};

export default ReviewPage;