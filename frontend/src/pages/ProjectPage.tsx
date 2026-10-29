import React, { useCallback, useEffect, useRef, useState } from "react";
import { Link, useParams } from "react-router-dom";
import api from "../api";
import { Review } from "../types";
import "./ProjectPage.css";

type ReviewMethod = "zip" | "single" | "paste" | "git";

const ProjectPage = () => {
    const { projectId } = useParams();

    const [activeMethod, setActiveMethod] = useState<ReviewMethod>("zip");
    const [zipFile, setZipFile] = useState<File | null>(null);
    const [singleFile, setSingleFile] = useState<File | null>(null);
    const [filename, setFilename] = useState("snippet.ts");
    const [language, setLanguage] = useState("typescript");
    const [code, setCode] = useState("// Paste code here\nconst greet = (name: string) => 'Hello ' + name;\nconsole.log(greet('world'));\n");
    const [url, setUrl] = useState("");
    const [branch, setBranch] = useState("main");
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    const [reviews, setReviews] = useState<Review[]>([]);
    const [uploadingZip, setUploadingZip] = useState(false);
    const [uploadingSingle, setUploadingSingle] = useState(false);
    const [pastingCode, setPastingCode] = useState(false);
    const [scanning, setScanning] = useState(false);
    const [loadingReviews, setLoadingReviews] = useState(true);

    const timerRef = useRef<number | null>(null);

    const loadReviews = useCallback(async () => {
        if (!projectId) return;

        try {
            const response = await api.get<Review[]>(`/api/projects/${projectId}/reviews`);
            setReviews(response.data);
            setError("");
        } catch (requestError) {
            console.error("Unable to load reviews:", requestError);
            setError("Unable to load review history.");
        } finally {
            setLoadingReviews(false);
        }
    }, [projectId]);

    useEffect(() => {
        loadReviews();
        timerRef.current = window.setInterval(() => {
            loadReviews();
        }, 3000);

        return () => {
            if (timerRef.current !== null) {
                window.clearInterval(timerRef.current);
            }
        };
    }, [loadReviews]);

    const handleZipUpload = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!zipFile || !projectId) return;

        try {
            setUploadingZip(true);
            setMessage("");
            setError("");

            const form = new FormData();
            form.append("file", zipFile, zipFile.name);

            await api.post(`/api/projects/${projectId}/reviews/upload`, form, {
                headers: { "Content-Type": "multipart/form-data" },
            });
            setMessage("ZIP upload accepted. Your code review is now running.");
            setZipFile(null);

            const zipInput = document.getElementById("projectZipFile") as HTMLInputElement | null;
            if (zipInput) zipInput.value = "";

            await loadReviews();
        } catch (requestError) {
            console.error("ZIP upload failed:", requestError);
            setError("Unable to upload the project ZIP. Please check the file and try again.");
        } finally {
            setUploadingZip(false);
        }
    };

    const handleSingleFileUpload = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!singleFile || !projectId) return;

        try {
            setUploadingSingle(true);
            setMessage("");
            setError("");

            const form = new FormData();
            form.append("file", singleFile, singleFile.name);

            await api.post(`/api/projects/${projectId}/reviews/upload-file`, form, {
                headers: { "Content-Type": "multipart/form-data" },
            });
            setMessage("Single-file review accepted. Your code review is now running.");
            setSingleFile(null);

            const input = document.getElementById("projectSingleFile") as HTMLInputElement | null;
            if (input) input.value = "";

            await loadReviews();
        } catch (requestError) {
            console.error("Single-file upload failed:", requestError);
            setError("Unable to upload the source file. Please verify the file type and try again.");
        } finally {
            setUploadingSingle(false);
        }
    };

    const handlePasteCode = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!code.trim() || !projectId) return;

        try {
            setPastingCode(true);
            setMessage("");
            setError("");

            await api.post(`/api/projects/${projectId}/reviews/paste`, {
                filename: filename.trim() || "pasted-code.ts",
                language: language.trim() || "typescript",
                code,
            });

            setMessage("Pasted code accepted. Your review is now running.");
            setCode("");
            await loadReviews();
        } catch (requestError) {
            console.error("Paste review failed:", requestError);
            setError("Unable to review the pasted code. Please check the code and try again.");
        } finally {
            setPastingCode(false);
        }
    };

    const handleGitReview = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!url.trim() || !projectId) return;

        try {
            setScanning(true);
            setMessage("");
            setError("");

            await api.post(`/api/projects/${projectId}/reviews/git`, {
                repositoryUrl: url.trim(),
                branch: branch.trim() || "main",
            });

            setMessage("Git repository accepted successfully. Code review is now running.");
            await loadReviews();
        } catch (requestError) {
            console.error("Git review failed:", requestError);
            setError("Unable to review the repository. Check the repository URL and branch and try again.");
        } finally {
            setScanning(false);
        }
    };

    const getStatusClass = (status: string) => {
        const normalizedStatus = status?.toLowerCase();
        if (normalizedStatus === "completed") return "review-status completed";
        if (normalizedStatus === "failed") return "review-status failed";
        if (normalizedStatus === "running" || normalizedStatus === "processing") return "review-status running";
        if (normalizedStatus === "pending" || normalizedStatus === "queued") return "review-status pending";
        return "review-status default";
    };

    const getSourceIcon = (sourceName: string) => {
        if (sourceName?.toLowerCase().includes("github")) return "bi-github";
        if (sourceName?.toLowerCase().endsWith(".zip")) return "bi-file-earmark-zip";
        return "bi-code-square";
    };

    const renderMethodContent = () => {
        if (activeMethod === "zip") {
            return (
                <form className="review-method-form" onSubmit={handleZipUpload}>
                    <div className="selected-method-heading">
                        <div className="selected-method-icon zip">
                            <i className="bi bi-cloud-arrow-up"></i>
                        </div>
                        <div>
                            <h2>Upload project ZIP</h2>
                            <p>Upload a full project archive to review the entire codebase in one pass.</p>
                        </div>
                    </div>

                    <label htmlFor="projectZipFile" className="project-field-label">Project file</label>
                    <div className={`project-file-dropzone ${zipFile ? "has-file" : ""}`}>
                        <div className="project-dropzone-icon">
                            <i className="bi bi-file-earmark-zip"></i>
                        </div>
                        <div className="project-dropzone-content">
                            <strong>{zipFile ? zipFile.name : "Choose your project ZIP"}</strong>
                            <span>{zipFile ? "Project is ready to upload" : "Select a .zip file from your computer"}</span>
                        </div>
                        <label htmlFor="projectZipFile" className="project-browse-button">Browse</label>
                        <input id="projectZipFile" name="file" type="file" accept=".zip" onChange={(event) => setZipFile(event.target.files?.[0] || null)} />
                    </div>

                    <div className="review-form-footer">
                        <div className="review-form-note">
                            <i className="bi bi-info-circle"></i>
                            ZIP files are extracted and scanned securely.
                        </div>
                        <button type="submit" className="project-primary-button" disabled={!zipFile || uploadingZip}>
                            {uploadingZip ? (
                                <>
                                    <span className="project-spinner small"></span>
                                    Uploading...
                                </>
                            ) : (
                                <>
                                    <i className="bi bi-cloud-arrow-up"></i>
                                    Upload and review
                                </>
                            )}
                        </button>
                    </div>
                </form>
            );
        }

        if (activeMethod === "single") {
            return (
                <form className="review-method-form" onSubmit={handleSingleFileUpload}>
                    <div className="selected-method-heading">
                        <div className="selected-method-icon file">
                            <i className="bi bi-file-earmark-code"></i>
                        </div>
                        <div>
                            <h2>Upload source file</h2>
                            <p>Review one uploaded source file such as JavaScript, TypeScript, Python, or Go.</p>
                        </div>
                    </div>

                    <label htmlFor="projectSingleFile" className="project-field-label">Source file</label>
                    <div className={`project-file-dropzone ${singleFile ? "has-file" : ""}`}>
                        <div className="project-dropzone-icon">
                            <i className="bi bi-file-earmark-text"></i>
                        </div>
                        <div className="project-dropzone-content">
                            <strong>{singleFile ? singleFile.name : "Choose a source file"}</strong>
                            <span>{singleFile ? "File is ready to review" : ".js, .ts, .tsx, .py, .go, .java, .sql, .json"}</span>
                        </div>
                        <label htmlFor="projectSingleFile" className="project-browse-button">Browse</label>
                        <input id="projectSingleFile" name="file" type="file" onChange={(event) => setSingleFile(event.target.files?.[0] || null)} />
                    </div>

                    <div className="review-form-footer">
                        <div className="review-form-note">
                            <i className="bi bi-shield-check"></i>
                            File review supports source and config files.
                        </div>
                        <button type="submit" className="project-primary-button" disabled={!singleFile || uploadingSingle}>
                            {uploadingSingle ? (
                                <>
                                    <span className="project-spinner small"></span>
                                    Uploading...
                                </>
                            ) : (
                                <>
                                    <i className="bi bi-file-earmark-code"></i>
                                    Review file
                                </>
                            )}
                        </button>
                    </div>
                </form>
            );
        }

        if (activeMethod === "paste") {
            return (
                <form className="review-method-form" onSubmit={handlePasteCode}>
                    <div className="selected-method-heading">
                        <div className="selected-method-icon paste">
                            <i className="bi bi-code-square"></i>
                        </div>
                        <div>
                            <h2>Paste code snippet</h2>
                            <p>Paste a function, module, or class to get a quick security and quality review.</p>
                        </div>
                    </div>

                    <div className="paste-code-options">
                        <div className="project-form-group">
                            <label htmlFor="pasteFilename" className="project-field-label">Filename</label>
                            <div className="project-input-wrapper">
                                <i className="bi bi-type"></i>
                                <input id="pasteFilename" className="project-input" value={filename} onChange={(event) => setFilename(event.target.value)} placeholder="snippet.ts" />
                            </div>
                        </div>

                        <div className="project-form-group">
                            <label htmlFor="pasteLanguage" className="project-field-label">Language</label>
                            <div className="project-input-wrapper">
                                <i className="bi bi-code"></i>
                                <select id="pasteLanguage" className="project-input project-select" value={language} onChange={(event) => setLanguage(event.target.value)}>
                                    <option value="typescript">TypeScript</option>
                                    <option value="javascript">JavaScript</option>
                                    <option value="python">Python</option>
                                    <option value="go">Go</option>
                                    <option value="java">Java</option>
                                    <option value="sql">SQL</option>
                                    <option value="json">JSON</option>
                                    <option value="bash">Bash</option>
                                </select>
                            </div>
                        </div>
                    </div>

                    <div className="project-form-group">
                        <label className="project-field-label">Code</label>
                        <div className="paste-code-editor">
                            <div className="paste-code-editor-header">
                                <div className="editor-dots">
                                    <span className="editor-dot red"></span>
                                    <span className="editor-dot yellow"></span>
                                    <span className="editor-dot green"></span>
                                </div>
                                <span>snippet</span>
                                <span>{filename || "untitled"}</span>
                            </div>
                            <textarea className="paste-code-textarea" value={code} onChange={(event) => setCode(event.target.value)} placeholder="Paste code here..." />
                        </div>
                    </div>

                    <div className="review-form-footer">
                        <div className="review-form-note">
                            <i className="bi bi-lightning-charge"></i>
                            Review runs immediately after you paste the code.
                        </div>
                        <button type="submit" className="project-primary-button" disabled={pastingCode || !code.trim()}>
                            {pastingCode ? (
                                <>
                                    <span className="project-spinner small"></span>
                                    Reviewing...
                                </>
                            ) : (
                                <>
                                    <i className="bi bi-code-square"></i>
                                    Review pasted code
                                </>
                            )}
                        </button>
                    </div>
                </form>
            );
        }

        return (
            <form className="review-method-form" onSubmit={handleGitReview}>
                <div className="selected-method-heading">
                    <div className="selected-method-icon git">
                        <i className="bi bi-github"></i>
                    </div>
                    <div>
                        <h2>Scan GitHub repository</h2>
                        <p>Connect a public Git repository using its HTTPS clone URL and review the selected branch.</p>
                    </div>
                </div>

                <div className="git-form-grid">
                    <div className="project-form-group">
                        <label htmlFor="repositoryUrl" className="project-field-label">Repository URL</label>
                        <div className="project-input-wrapper">
                            <i className="bi bi-git"></i>
                            <input id="repositoryUrl" type="url" className="project-input" placeholder="https://github.com/user/repo.git" value={url} onChange={(event) => setUrl(event.target.value)} required />
                        </div>
                    </div>

                    <div className="project-form-group">
                        <label htmlFor="repositoryBranch" className="project-field-label">Branch</label>
                        <div className="project-input-wrapper">
                            <i className="bi bi-diagram-2"></i>
                            <input id="repositoryBranch" className="project-input" placeholder="main" value={branch} onChange={(event) => setBranch(event.target.value)} />
                        </div>
                    </div>
                </div>

                <div className="review-form-footer">
                    <div className="review-form-note">
                        <i className="bi bi-git"></i>
                        Repository review checks the selected branch for security and quality issues.
                    </div>
                    <button type="submit" className="project-primary-button" disabled={scanning || !url.trim()}>
                        {scanning ? (
                            <>
                                <span className="project-spinner small"></span>
                                Starting review...
                            </>
                        ) : (
                            <>
                                <i className="bi bi-github"></i>
                                Clone and review
                            </>
                        )}
                    </button>
                </div>
            </form>
        );
    };

    return (
        <div className="project-review-page">
            <div className="project-review-container">
                <div className="project-review-breadcrumb">
                    <Link to="/dashboard">
                        <i className="bi bi-arrow-left"></i>
                        Back to projects
                    </Link>
                </div>

                <div className="project-review-header">
                    <div>
                        <div className="project-review-label">
                            <i className="bi bi-stars"></i>
                            CODE REVIEW
                        </div>
                        <h1>Project reviews</h1>
                        <p>
                            Upload your source code, paste a snippet, or connect a Git repository to identify bugs,
                            security risks, performance issues, and code-quality improvements.
                        </p>
                    </div>

                    <Link to="/dashboard" className="project-dashboard-button">
                        <i className="bi bi-grid"></i>
                        Dashboard
                    </Link>
                </div>

                {message && (
                    <div className="project-message success">
                        <div className="project-message-icon">
                            <i className="bi bi-check-circle-fill"></i>
                        </div>
                        <div>
                            <strong>Review started</strong>
                            <span>{message}</span>
                        </div>
                        <button type="button" onClick={() => setMessage("")} aria-label="Close message">
                            <i className="bi bi-x-lg"></i>
                        </button>
                    </div>
                )}

                {error && (
                    <div className="project-message error">
                        <div className="project-message-icon">
                            <i className="bi bi-exclamation-circle-fill"></i>
                        </div>
                        <div>
                            <strong>Something went wrong</strong>
                            <span>{error}</span>
                        </div>
                        <button type="button" onClick={() => setError("")} aria-label="Close error">
                            <i className="bi bi-x-lg"></i>
                        </button>
                    </div>
                )}

                <div className="review-selection-section">
                    <div className="review-selection-heading">
                        <span className="review-section-label">REVIEW METHOD</span>
                        <h2>Choose how to review your code</h2>
                        <p>Upload a project archive, a single file, paste code, or connect a Git repository.</p>
                    </div>

                    <div className="review-method-tabs">
                        {[
                            { id: "zip", label: "ZIP File", desc: "Full project upload", icon: "bi-cloud-arrow-up" },
                            { id: "single", label: "Single File", desc: "One source file", icon: "bi-file-earmark-code" },
                            { id: "paste", label: "Paste Code", desc: "Snippet review", icon: "bi-code-square" },
                            { id: "git", label: "Git Repo", desc: "Clone and scan", icon: "bi-github" },
                        ].map((method) => (
                            <button
                                key={method.id}
                                type="button"
                                className={`review-method-tab ${activeMethod === method.id ? "active" : ""}`}
                                onClick={() => setActiveMethod(method.id as ReviewMethod)}
                            >
                                <div className="review-method-tab-icon">
                                    <i className={method.icon}></i>
                                </div>
                                <div className="review-method-tab-content">
                                    <strong>{method.label}</strong>
                                    <span>{method.desc}</span>
                                </div>
                                <div className="review-method-selected">
                                    <i className="bi bi-check-circle-fill"></i>
                                </div>
                            </button>
                        ))}
                    </div>

                    <div className="review-method-content">{renderMethodContent()}</div>
                </div>

                <section className="review-history-section">
                    <div className="review-history-header">
                        <div>
                            <div className="review-history-title-row">
                                <div className="review-history-icon">
                                    <i className="bi bi-clock-history"></i>
                                </div>
                                <div>
                                    <h2>Review history</h2>
                                    <p>Monitor previous and currently running code reviews.</p>
                                </div>
                            </div>
                        </div>
                        <button type="button" className="review-refresh-button" onClick={loadReviews}>
                            <i className="bi bi-arrow-clockwise"></i>
                            Refresh
                        </button>
                    </div>

                    {loadingReviews ? (
                        <div className="review-loading">
                            <div className="project-spinner large"></div>
                            <span>Loading review history...</span>
                        </div>
                    ) : reviews.length === 0 ? (
                        <div className="review-empty-state">
                            <div className="review-empty-icon">
                                <i className="bi bi-code-square"></i>
                            </div>
                            <h3>No reviews yet</h3>
                            <p>Upload a project ZIP, source file, or connect a Git repository to start your first code review.</p>
                        </div>
                    ) : (
                        <div className="review-table-wrapper">
                            <table className="review-table">
                                <thead>
                                    <tr>
                                        <th>Source</th>
                                        <th>Status</th>
                                        <th>Files</th>
                                        <th>Critical</th>
                                        <th>High</th>
                                        <th>Medium</th>
                                        <th>Low</th>
                                        <th>Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {reviews.map((review) => (
                                        <tr key={review.id}>
                                            <td>
                                                <div className="review-source-name">
                                                    <div className="review-table-source-icon">
                                                        <i className={`bi ${getSourceIcon(review.sourceName)}`}></i>
                                                    </div>
                                                    <span>{review.sourceName}</span>
                                                </div>
                                            </td>
                                            <td>
                                                <span className={getStatusClass(review.status)}>
                                                    <span className="review-status-dot"></span>
                                                    {review.status}
                                                </span>
                                            </td>
                                            <td><span className="review-file-count">{review.filesScanned}</span></td>
                                            <td><span className={`severity-value critical ${review.criticalCount > 0 ? "has-issues" : ""}`}>{review.criticalCount}</span></td>
                                            <td><span className={`severity-value high ${review.highCount > 0 ? "has-issues" : ""}`}>{review.highCount}</span></td>
                                            <td><span className={`severity-value medium ${review.mediumCount > 0 ? "has-issues" : ""}`}>{review.mediumCount}</span></td>
                                            <td><span className={`severity-value low ${review.lowCount > 0 ? "has-issues" : ""}`}>{review.lowCount}</span></td>
                                            <td>
                                                <Link to={`/reviews/${review.id}`} className="review-open-button">
                                                    View report
                                                    <i className="bi bi-arrow-right"></i>
                                                </Link>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    )}
                </section>
            </div>
        </div>
    );
};

export default ProjectPage;
