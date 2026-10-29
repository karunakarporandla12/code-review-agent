// import React, { useCallback, useEffect, useState } from "react";
// import api from "../api";

// type FrontendTechnology = "react" | "angular";
// type BackendTechnology = "java" | "golang";
// type DatabaseTechnology = "postgresql" | "mysql" | "mongodb" | "none";

// interface GeneratedProject {
//   id: string;
//   projectName: string;
//   requirements: string;
//   frontendTechnology: FrontendTechnology;
//   backendTechnology: BackendTechnology;
//   databaseTechnology: DatabaseTechnology;
//   gitUrl: string;
//   gitBranch: string;
//   rootPath: string;
//   files: string[];
//   createdAt?: string;
//   updatedAt?: string;
// }

// interface GenerateProjectResponse {
//   id: string;
//   root: string;
//   files: string[];
// }

// const FeatureGeneratorPage: React.FC = () => {
//   const [projectName, setProjectName] = useState("");
//   const [requirements, setRequirements] = useState("");
//   const [frontendTechnology, setFrontendTechnology] =
//     useState<FrontendTechnology>("react");
//   const [backendTechnology, setBackendTechnology] =
//     useState<BackendTechnology>("golang");
//   const [databaseTechnology, setDatabaseTechnology] =
//     useState<DatabaseTechnology>("postgresql");
//   const [gitUrl, setGitUrl] = useState("");
//   const [gitBranch, setGitBranch] = useState("main");

//   const [loading, setLoading] = useState(false);
//   const [loadingProjects, setLoadingProjects] = useState(false);
//   const [loadingFile, setLoadingFile] = useState(false);
//   const [savingFile, setSavingFile] = useState(false);
//   const [syncingGit, setSyncingGit] = useState(false);

//   const [generatedFiles, setGeneratedFiles] = useState<string[]>([]);
//   const [rootPath, setRootPath] = useState("");
//   const [savedProjects, setSavedProjects] = useState<GeneratedProject[]>([]);
//   const [selectedProjectId, setSelectedProjectId] = useState("");
//   const [selectedFile, setSelectedFile] = useState("");
//   const [fileContent, setFileContent] = useState("");

//   const [error, setError] = useState("");
//   const [success, setSuccess] = useState("");

//   const clearMessages = () => {
//     setError("");
//     setSuccess("");
//   };

//   const resetFileEditor = () => {
//     setSelectedFile("");
//     setFileContent("");
//   };

//   const loadSavedProjects = useCallback(async () => {
//     setLoadingProjects(true);

//     try {
//       const response = await api.get<GeneratedProject[]>(
//         "/api/projects/generated"
//       );

//       const projects = Array.isArray(response.data) ? response.data : [];
//       setSavedProjects(projects);

//       setSelectedProjectId((currentProjectId) => {
//         if (currentProjectId) {
//           return currentProjectId;
//         }

//         return projects.length > 0 ? projects[0].id : "";
//       });
//     } catch (requestError) {
//       console.error("Failed to load generated projects:", requestError);
//       setError("Unable to load the saved generated projects.");
//     } finally {
//       setLoadingProjects(false);
//     }
//   }, []);

//   const loadProject = useCallback(async (projectId: string) => {
//     if (!projectId) {
//       return;
//     }

//     clearMessages();
//     resetFileEditor();

//     try {
//       const response = await api.get<GeneratedProject>(
//         `/api/projects/generated/${projectId}`
//       );

//       const project = response.data;

//       setProjectName(project.projectName || "");
//       setRequirements(project.requirements || "");
//       setFrontendTechnology(project.frontendTechnology || "react");
//       setBackendTechnology(project.backendTechnology || "golang");
//       setDatabaseTechnology(project.databaseTechnology || "postgresql");
//       setGitUrl(project.gitUrl || "");
//       setGitBranch(project.gitBranch || "main");
//       setRootPath(project.rootPath || "");
//       setGeneratedFiles(project.files || []);
//     } catch (requestError) {
//       console.error("Failed to load generated project:", requestError);
//       setError("Unable to load the selected generated project.");
//     }
//   }, []);

//   useEffect(() => {
//     loadSavedProjects();
//   }, [loadSavedProjects]);

//   useEffect(() => {
//     if (selectedProjectId) {
//       loadProject(selectedProjectId);
//     }
//   }, [selectedProjectId, loadProject]);

//   const validateGenerateForm = (): boolean => {
//     if (!projectName.trim()) {
//       setError("Please enter a project name.");
//       return false;
//     }

//     if (!requirements.trim()) {
//       setError("Please enter the project requirements.");
//       return false;
//     }

//     return true;
//   };

//   const handleGenerate = async () => {
//     clearMessages();

//     if (!validateGenerateForm()) {
//       return;
//     }

//     setLoading(true);
//     setGeneratedFiles([]);
//     setRootPath("");
//     resetFileEditor();

//     try {
//       const response = await api.post<GenerateProjectResponse>(
//         "/api/projects/generate",
//         {
//           projectName: projectName.trim(),
//           requirements: requirements.trim(),
//           frontendTechnology,
//           backendTechnology,
//           databaseTechnology,
//           gitUrl: gitUrl.trim(),
//           gitBranch: gitBranch.trim() || "main",
//         }
//       );

//       setSelectedProjectId(response.data.id);
//       setRootPath(response.data.root || "");
//       setGeneratedFiles(response.data.files || []);
//       setSuccess("The project was generated successfully.");

//       await loadSavedProjects();
//     } catch (requestError) {
//       console.error("Project generation failed:", requestError);
//       setError(
//         "Unable to generate the project. Check the backend logs and database connection."
//       );
//     } finally {
//       setLoading(false);
//     }
//   };

//   const openFile = async (filePath: string) => {
//     if (!selectedProjectId || !filePath) {
//       return;
//     }

//     clearMessages();
//     setLoadingFile(true);

//     try {
//       const response = await api.get<{ content: string }>(
//         `/api/projects/generated/${selectedProjectId}/file`,
//         {
//           params: {
//             path: filePath,
//           },
//         }
//       );

//       setSelectedFile(filePath);
//       setFileContent(response.data.content || "");
//     } catch (requestError) {
//       console.error("Failed to open generated file:", requestError);
//       setError("Unable to open the selected file.");
//     } finally {
//       setLoadingFile(false);
//     }
//   };

//   const saveFile = async () => {
//     if (!selectedProjectId || !selectedFile) {
//       return;
//     }

//     clearMessages();
//     setSavingFile(true);

//     try {
//       await api.put(
//         `/api/projects/generated/${selectedProjectId}/file`,
//         {
//           path: selectedFile,
//           content: fileContent,
//         }
//       );

//       setSuccess("The file was saved successfully.");
//     } catch (requestError) {
//       console.error("Failed to save generated file:", requestError);
//       setError("Unable to save the selected file.");
//     } finally {
//       setSavingFile(false);
//     }
//   };

//   const syncGit = async () => {
//     if (!selectedProjectId) {
//       setError("Please select or generate a project.");
//       return;
//     }

//     if (!gitUrl.trim()) {
//       setError("Please enter a Git remote URL.");
//       return;
//     }

//     clearMessages();
//     setSyncingGit(true);

//     try {
//       await api.post(
//         `/api/projects/generated/${selectedProjectId}/sync-git`,
//         {
//           gitUrl: gitUrl.trim(),
//           gitBranch: gitBranch.trim() || "main",
//         }
//       );

//       setSuccess("The generated project was pushed to Git successfully.");
//       await loadSavedProjects();
//     } catch (requestError) {
//       console.error("Git synchronization failed:", requestError);
//       setError(
//         "Unable to push the project to Git. Check the repository URL and backend Git credentials."
//       );
//     } finally {
//       setSyncingGit(false);
//     }
//   };

//   const startNewProject = () => {
//     clearMessages();
//     setSelectedProjectId("");
//     setProjectName("");
//     setRequirements("");
//     setFrontendTechnology("react");
//     setBackendTechnology("golang");
//     setDatabaseTechnology("postgresql");
//     setGitUrl("");
//     setGitBranch("main");
//     setGeneratedFiles([]);
//     setRootPath("");
//     resetFileEditor();
//   };

//   return (
//     <div className="container py-4">
//       <div className="card border-0 shadow-sm">
//         <div className="card-body p-4">
//           <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
//             <div>
//               <h2 className="mb-2">Generate Application From Requirements</h2>
//               <p className="text-muted mb-0">
//                 Select the required technology stack and describe the
//                 application. A separate project folder will be created for
//                 every generated application.
//               </p>
//             </div>

//             <button
//               type="button"
//               className="btn btn-outline-primary"
//               onClick={startNewProject}
//             >
//               New Project
//             </button>
//           </div>

//           {error && (
//             <div className="alert alert-danger" role="alert">
//               {error}
//             </div>
//           )}

//           {success && (
//             <div className="alert alert-success" role="alert">
//               {success}
//             </div>
//           )}

//           <div className="row g-4">
//             <div className="col-lg-7">
//               <div className="mb-3">
//                 <label htmlFor="projectName" className="form-label">
//                   Project Name
//                 </label>

//                 <input
//                   id="projectName"
//                   type="text"
//                   className="form-control"
//                   value={projectName}
//                   onChange={(event) => setProjectName(event.target.value)}
//                   placeholder="For example: fashion-shop"
//                   disabled={loading}
//                 />
//               </div>

//               <div className="row g-3 mb-3">
//                 <div className="col-md-4">
//                   <label htmlFor="frontendTechnology" className="form-label">
//                     Frontend
//                   </label>

//                   <select
//                     id="frontendTechnology"
//                     className="form-select"
//                     value={frontendTechnology}
//                     onChange={(event) =>
//                       setFrontendTechnology(
//                         event.target.value as FrontendTechnology
//                       )
//                     }
//                     disabled={loading}
//                   >
//                     <option value="react">React</option>
//                     <option value="angular">Angular</option>
//                   </select>
//                 </div>

//                 <div className="col-md-4">
//                   <label htmlFor="backendTechnology" className="form-label">
//                     Backend
//                   </label>

//                   <select
//                     id="backendTechnology"
//                     className="form-select"
//                     value={backendTechnology}
//                     onChange={(event) =>
//                       setBackendTechnology(
//                         event.target.value as BackendTechnology
//                       )
//                     }
//                     disabled={loading}
//                   >
//                     <option value="golang">Go</option>
//                     <option value="java">Java Spring Boot</option>
//                   </select>
//                 </div>

//                 <div className="col-md-4">
//                   <label htmlFor="databaseTechnology" className="form-label">
//                     Database
//                   </label>

//                   <select
//                     id="databaseTechnology"
//                     className="form-select"
//                     value={databaseTechnology}
//                     onChange={(event) =>
//                       setDatabaseTechnology(
//                         event.target.value as DatabaseTechnology
//                       )
//                     }
//                     disabled={loading}
//                   >
//                     <option value="postgresql">PostgreSQL</option>
//                     <option value="mysql">MySQL</option>
//                     <option value="mongodb">MongoDB</option>
//                     <option value="none">No Database</option>
//                   </select>
//                 </div>
//               </div>

//               <div className="mb-3">
//                 <label htmlFor="requirements" className="form-label">
//                   Application Requirements
//                 </label>

//                 <textarea
//                   id="requirements"
//                   className="form-control"
//                   rows={10}
//                   value={requirements}
//                   onChange={(event) => setRequirements(event.target.value)}
//                   placeholder="Example: Create a clothing store with user registration, login, product catalog, category filters, shopping cart, checkout, order history, admin product management and sales dashboard."
//                   disabled={loading}
//                 />
//               </div>

//               <div className="mb-3">
//                 <label htmlFor="gitUrl" className="form-label">
//                   Git Remote URL
//                 </label>

//                 <input
//                   id="gitUrl"
//                   type="text"
//                   className="form-control"
//                   value={gitUrl}
//                   onChange={(event) => setGitUrl(event.target.value)}
//                   placeholder="Git repository URL"
//                   disabled={loading}
//                 />
//               </div>

//               <div className="mb-4">
//                 <label htmlFor="gitBranch" className="form-label">
//                   Git Branch
//                 </label>

//                 <input
//                   id="gitBranch"
//                   type="text"
//                   className="form-control"
//                   value={gitBranch}
//                   onChange={(event) => setGitBranch(event.target.value)}
//                   placeholder="main"
//                   disabled={loading}
//                 />
//               </div>

//               <button
//                 type="button"
//                 className="btn btn-primary"
//                 onClick={handleGenerate}
//                 disabled={
//                   loading ||
//                   !projectName.trim() ||
//                   !requirements.trim()
//                 }
//               >
//                 {loading ? "Generating Project..." : "Generate Project"}
//               </button>
//             </div>

//             <div className="col-lg-5">
//               <div className="border rounded p-3 h-100">
//                 <h4 className="mb-3">Saved Projects</h4>

//                 {loadingProjects && (
//                   <div className="text-muted">Loading projects...</div>
//                 )}

//                 {!loadingProjects && savedProjects.length === 0 && (
//                   <div className="text-muted">
//                     No generated projects are available.
//                   </div>
//                 )}

//                 {!loadingProjects && savedProjects.length > 0 && (
//                   <div className="list-group">
//                     {savedProjects.map((project) => (
//                       <button
//                         key={project.id}
//                         type="button"
//                         className={`list-group-item list-group-item-action ${
//                           selectedProjectId === project.id ? "active" : ""
//                         }`}
//                         onClick={() => setSelectedProjectId(project.id)}
//                       >
//                         <div className="fw-semibold">
//                           {project.projectName}
//                         </div>

//                         <small>
//                           {project.frontendTechnology} +{" "}
//                           {project.backendTechnology} +{" "}
//                           {project.databaseTechnology}
//                         </small>

//                         <div>
//                           <small>
//                             {project.gitUrl ? "Git configured" : "Local only"}
//                           </small>
//                         </div>
//                       </button>
//                     ))}
//                   </div>
//                 )}
//               </div>
//             </div>
//           </div>

//           {rootPath && (
//             <div className="alert alert-info mt-4 mb-0">
//               Generated project folder: <strong>{rootPath}</strong>
//             </div>
//           )}

//           {generatedFiles.length > 0 && (
//             <div className="mt-4">
//               <h4 className="mb-3">Generated Files</h4>

//               <div className="list-group">
//                 {generatedFiles.map((filePath) => (
//                   <button
//                     key={filePath}
//                     type="button"
//                     className="list-group-item list-group-item-action text-start"
//                     onClick={() => openFile(filePath)}
//                     disabled={loadingFile}
//                   >
//                     {filePath}
//                   </button>
//                 ))}
//               </div>
//             </div>
//           )}

//           {selectedFile && (
//             <div className="mt-4 border rounded p-3">
//               <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-3">
//                 <strong>{selectedFile}</strong>

//                 <div className="d-flex flex-wrap gap-2">
//                   <button
//                     type="button"
//                     className="btn btn-sm btn-outline-primary"
//                     onClick={saveFile}
//                     disabled={savingFile}
//                   >
//                     {savingFile ? "Saving..." : "Save File"}
//                   </button>

//                   <button
//                     type="button"
//                     className="btn btn-sm btn-outline-success"
//                     onClick={syncGit}
//                     disabled={syncingGit || !gitUrl.trim()}
//                   >
//                     {syncingGit ? "Pushing..." : "Push To Git"}
//                   </button>
//                 </div>
//               </div>

//               <textarea
//                 className="form-control font-monospace"
//                 rows={22}
//                 value={fileContent}
//                 onChange={(event) => setFileContent(event.target.value)}
//                 spellCheck={false}
//               />
//             </div>
//           )}
//         </div>
//       </div>
//     </div>
//   );
// };

// export default FeatureGeneratorPage;
import React, {
    useCallback,
    useEffect,
    useState,
} from "react";
import {
    Link,
    useNavigate,
    useParams,
} from "react-router-dom";
import api from "../api";
import "./FeatureGeneratorPage.css";

type FrontendTechnology =
    | "react"
    | "angular";

type BackendTechnology =
    | "java"
    | "golang";

type DatabaseTechnology =
    | "postgresql"
    | "mysql"
    | "mongodb"
    | "none";

interface GeneratedProject {
    id: string;
    projectName: string;
    requirements: string;
    frontendTechnology: FrontendTechnology;
    backendTechnology: BackendTechnology;
    databaseTechnology: DatabaseTechnology;
    gitUrl: string;
    gitBranch: string;
    rootPath: string;
    files: string[];
    createdAt?: string;
    updatedAt?: string;
}

interface FileContentResponse {
    content: string;
}

const FeatureGeneratorPage: React.FC = () => {
    const navigate = useNavigate();

    const { projectId } = useParams<{
        projectId: string;
    }>();

    const [project, setProject] =
        useState<GeneratedProject | null>(null);

    const [selectedFile, setSelectedFile] =
        useState("");

    const [fileContent, setFileContent] =
        useState("");

    const [loadingProject, setLoadingProject] =
        useState(true);

    const [loadingFile, setLoadingFile] =
        useState(false);

    const [savingFile, setSavingFile] =
        useState(false);

    const [syncingGit, setSyncingGit] =
        useState(false);

    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const clearMessages = () => {
        setError("");
        setSuccess("");
    };

    const loadProject = useCallback(async () => {
        if (!projectId) {
            setError(
                "Generated project ID is missing."
            );

            setLoadingProject(false);
            return;
        }

        try {
            setLoadingProject(true);
            setError("");
            setSuccess("");
            setSelectedFile("");
            setFileContent("");

            const response =
                await api.get<GeneratedProject>(
                    `/api/projects/generated/${projectId}`
                );

            const projectData = response.data;

            setProject({
                ...projectData,
                files: Array.isArray(
                    projectData.files
                )
                    ? projectData.files
                    : [],
            });
        } catch (requestError) {
            console.error(
                "Failed to load generated project:",
                requestError
            );

            setProject(null);

            setError(
                "Unable to load the selected generated project."
            );
        } finally {
            setLoadingProject(false);
        }
    }, [projectId]);

    useEffect(() => {
        loadProject();
    }, [loadProject]);

    const openFile = async (
        filePath: string
    ) => {
        if (!projectId || !filePath) {
            return;
        }

        try {
            clearMessages();
            setLoadingFile(true);

            const response =
                await api.get<FileContentResponse>(
                    `/api/projects/generated/${projectId}/file`,
                    {
                        params: {
                            path: filePath,
                        },
                    }
                );

            setSelectedFile(filePath);

            setFileContent(
                response.data.content || ""
            );
        } catch (requestError) {
            console.error(
                "Failed to open generated file:",
                requestError
            );

            setError(
                "Unable to open the selected file."
            );
        } finally {
            setLoadingFile(false);
        }
    };

    const saveFile = async () => {
        if (
            !projectId ||
            !selectedFile
        ) {
            return;
        }

        try {
            clearMessages();
            setSavingFile(true);

            await api.put(
                `/api/projects/generated/${projectId}/file`,
                {
                    path: selectedFile,
                    content: fileContent,
                }
            );

            setSuccess(
                "File saved successfully."
            );
        } catch (requestError) {
            console.error(
                "Failed to save generated file:",
                requestError
            );

            setError(
                "Unable to save the selected file."
            );
        } finally {
            setSavingFile(false);
        }
    };

    const syncGit = async () => {
        if (!projectId || !project) {
            return;
        }

        if (!project.gitUrl?.trim()) {
            setError(
                "This project does not have a Git remote URL."
            );
            return;
        }

        try {
            clearMessages();
            setSyncingGit(true);

            await api.post(
                `/api/projects/generated/${projectId}/sync-git`,
                {
                    gitUrl: project.gitUrl.trim(),
                    gitBranch:
                        project.gitBranch?.trim() ||
                        "main",
                }
            );

            setSuccess(
                "Project pushed to Git successfully."
            );
        } catch (requestError) {
            console.error(
                "Git synchronization failed:",
                requestError
            );

            setError(
                "Unable to push the project to Git. Check the Git URL and backend credentials."
            );
        } finally {
            setSyncingGit(false);
        }
    };

    if (loadingProject) {
        return (
            <div className="generated-project-page">
                <div className="generated-project-loading">
                    <div
                        className="spinner-border"
                        role="status"
                    >
                        <span className="visually-hidden">
                            Loading project
                        </span>
                    </div>

                    <span>
                        Loading generated project...
                    </span>
                </div>
            </div>
        );
    }

    if (!project) {
        return (
            <div className="generated-project-page">
                <div className="generated-project-error-state">
                    <div className="error-state-icon">
                        <i className="bi bi-exclamation-circle"></i>
                    </div>

                    <h2>
                        Project could not be loaded
                    </h2>

                    <p>
                        {error ||
                            "The selected generated project is unavailable."}
                    </p>

                    <button
                        type="button"
                        className="generated-back-button"
                        onClick={() =>
                            navigate("/dashboard")
                        }
                    >
                        <i className="bi bi-arrow-left"></i>
                        <span>
                            Return To Dashboard
                        </span>
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="generated-project-page">
            <div className="generated-project-topbar">
                <Link
                    to="/dashboard"
                    className="generated-back-link"
                >
                    <i className="bi bi-arrow-left"></i>
                    <span>Dashboard</span>
                </Link>

                <div className="generated-project-actions">
                    <button
                        type="button"
                        className="generated-refresh-button"
                        onClick={loadProject}
                    >
                        <i className="bi bi-arrow-clockwise"></i>
                        <span>Refresh</span>
                    </button>

                    <button
                        type="button"
                        className="generated-git-button"
                        onClick={syncGit}
                        disabled={
                            syncingGit ||
                            !project.gitUrl?.trim()
                        }
                    >
                        {syncingGit ? (
                            <>
                                <span
                                    className="spinner-border spinner-border-sm"
                                    aria-hidden="true"
                                ></span>

                                <span>Pushing</span>
                            </>
                        ) : (
                            <>
                                <i className="bi bi-git"></i>
                                <span>Push To Git</span>
                            </>
                        )}
                    </button>
                </div>
            </div>

            {error && (
                <div className="generated-message generated-error-message">
                    <i className="bi bi-exclamation-circle"></i>

                    <span>{error}</span>

                    <button
                        type="button"
                        onClick={() => setError("")}
                        aria-label="Close error"
                    >
                        <i className="bi bi-x-lg"></i>
                    </button>
                </div>
            )}

            {success && (
                <div className="generated-message generated-success-message">
                    <i className="bi bi-check-circle"></i>

                    <span>{success}</span>

                    <button
                        type="button"
                        onClick={() =>
                            setSuccess("")
                        }
                        aria-label="Close success message"
                    >
                        <i className="bi bi-x-lg"></i>
                    </button>
                </div>
            )}

            <section className="generated-project-header">
                <div className="generated-project-title-area">
                    <div className="generated-project-icon">
                        <i className="bi bi-code-square"></i>
                    </div>

                    <div>
                        <span className="generated-project-label">
                            Generated Application
                        </span>

                        <h1>
                            {project.projectName}
                        </h1>

                        <p>
                            {project.requirements ||
                                "No requirements available."}
                        </p>
                    </div>
                </div>

                <div className="generated-project-stack">
                    <div>
                        <span>Frontend</span>
                        <strong>
                            {
                                project.frontendTechnology
                            }
                        </strong>
                    </div>

                    <div>
                        <span>Backend</span>
                        <strong>
                            {
                                project.backendTechnology
                            }
                        </strong>
                    </div>

                    <div>
                        <span>Database</span>
                        <strong>
                            {
                                project.databaseTechnology
                            }
                        </strong>
                    </div>
                </div>
            </section>

            <section className="generated-project-details">
                <div className="generated-detail">
                    <span>Project ID</span>
                    <strong>{project.id}</strong>
                </div>

                <div className="generated-detail">
                    <span>Generated Folder</span>
                    <strong>
                        {project.rootPath ||
                            "Not available"}
                    </strong>
                </div>

                <div className="generated-detail">
                    <span>Git Repository</span>
                    <strong>
                        {project.gitUrl ||
                            "Local only"}
                    </strong>
                </div>

                <div className="generated-detail">
                    <span>Git Branch</span>
                    <strong>
                        {project.gitBranch ||
                            "main"}
                    </strong>
                </div>
            </section>

            <div className="generated-workspace">
                <aside className="generated-files-panel">
                    <div className="generated-panel-header">
                        <div>
                            <h2>Project Files</h2>

                            <p>
                                {project.files.length}{" "}
                                {project.files.length === 1
                                    ? "file"
                                    : "files"}
                            </p>
                        </div>

                        <i className="bi bi-folder2-open"></i>
                    </div>

                    {project.files.length === 0 ? (
                        <div className="generated-files-empty">
                            <i className="bi bi-file-earmark-x"></i>

                            <span>
                                No generated files were
                                returned for this project.
                            </span>
                        </div>
                    ) : (
                        <div className="generated-files-list">
                            {project.files.map(
                                (filePath) => (
                                    <button
                                        key={filePath}
                                        type="button"
                                        className={
                                            selectedFile ===
                                            filePath
                                                ? "generated-file-button active"
                                                : "generated-file-button"
                                        }
                                        onClick={() =>
                                            openFile(
                                                filePath
                                            )
                                        }
                                        disabled={
                                            loadingFile
                                        }
                                    >
                                        <i className="bi bi-file-earmark-code"></i>

                                        <span>
                                            {filePath}
                                        </span>
                                    </button>
                                )
                            )}
                        </div>
                    )}
                </aside>

                <section className="generated-editor-panel">
                    {!selectedFile ? (
                        <div className="generated-editor-empty">
                            <div className="editor-empty-icon">
                                <i className="bi bi-file-earmark-code"></i>
                            </div>

                            <h2>Select a project file</h2>

                            <p>
                                Choose a file from the
                                project file list to view
                                and edit the content.
                            </p>
                        </div>
                    ) : (
                        <>
                            <div className="generated-editor-header">
                                <div>
                                    <span>
                                        Selected File
                                    </span>

                                    <strong>
                                        {selectedFile}
                                    </strong>
                                </div>

                                <button
                                    type="button"
                                    className="generated-save-button"
                                    onClick={saveFile}
                                    disabled={savingFile}
                                >
                                    {savingFile ? (
                                        <>
                                            <span
                                                className="spinner-border spinner-border-sm"
                                                aria-hidden="true"
                                            ></span>

                                            <span>
                                                Saving
                                            </span>
                                        </>
                                    ) : (
                                        <>
                                            <i className="bi bi-floppy"></i>

                                            <span>
                                                Save File
                                            </span>
                                        </>
                                    )}
                                </button>
                            </div>

                            {loadingFile ? (
                                <div className="generated-file-loading">
                                    <div
                                        className="spinner-border"
                                        role="status"
                                    >
                                        <span className="visually-hidden">
                                            Loading file
                                        </span>
                                    </div>

                                    <span>
                                        Loading file...
                                    </span>
                                </div>
                            ) : (
                                <textarea
                                    className="generated-code-editor"
                                    value={fileContent}
                                    onChange={(event) =>
                                        setFileContent(
                                            event.target
                                                .value
                                        )
                                    }
                                    spellCheck={false}
                                    aria-label="Generated file content"
                                />
                            )}
                        </>
                    )}
                </section>
            </div>
        </div>
    );
};

export default FeatureGeneratorPage;