// import React, { useEffect, useState } from "react";
// import { Link } from "react-router-dom";
// import api from "../api";
// import "./DashboardPage.css";

// const DashboardPage = () => {
//     const [items, setItems] = useState([]);
//     const [generatedApps, setGeneratedApps] = useState([]);
//     const [name, setName] = useState("");
//     const [description, setDescription] = useState("");
//     const [loading, setLoading] = useState(true);
//     const [creating, setCreating] = useState(false);
//     const [error, setError] = useState("");

//     const loadProjects = async () => {
//         try {
//             setLoading(true);
//             setError("");

//             const response = await api.get("/api/projects");
//             setItems(response.data);

//             const generatedResponse = await api.get("/api/projects/generated");
//             setGeneratedApps(generatedResponse.data || []);
//         } catch (err) {
//             console.error("Failed to load projects:", err);
//             setError("Unable to load projects. Please try again.");
//         } finally {
//             setLoading(false);
//         }
//     };

//     useEffect(() => {
//         loadProjects();
//     }, []);

//     const createProject = async (event) => {
//         event.preventDefault();

//         const trimmedName = name.trim();
//         const trimmedDescription = description.trim();

//         if (!trimmedName) {
//             setError("Project name is required.");
//             return;
//         }

//         try {
//             setCreating(true);
//             setError("");

//             await api.post("/api/projects", {
//                 name: trimmedName,
//                 description: trimmedDescription,
//             });

//             setName("");
//             setDescription("");

//             await loadProjects();
//         } catch (err) {
//             console.error("Failed to create project:", err);
//             setError("Unable to create project. Please try again.");
//         } finally {
//             setCreating(false);
//         }
//     };

//     const formatCreatedDate = (createdAt) => {
//         if (!createdAt) {
//             return "Unknown";
//         }

//         const date = new Date(createdAt);

//         if (Number.isNaN(date.getTime())) {
//             return "Unknown";
//         }

//         return date.toLocaleDateString(undefined, {
//             day: "2-digit",
//             month: "short",
//             year: "numeric",
//         });
//     };

//     return (
//         <div className="dashboard-page">
//             <div className="dashboard-header">
//                 <div>
//                     <div className="dashboard-label">
//                         <i className="bi bi-grid"></i>
//                         <span>Dashboard</span>
//                     </div>

//                     <h1>Projects</h1>

//                     <p>
//                         Create and manage your projects, upload source code,
//                         or start a review from a Git repository.
//                     </p>
//                 </div>

//                 <div className="dashboard-summary">
//                     <div className="summary-icon">
//                         <i className="bi bi-folder2-open"></i>
//                     </div>

//                     <div>
//                         <span>Total Projects</span>
//                         <strong>{items.length}</strong>
//                     </div>
//                 </div>
//             </div>

//             {error && (
//                 <div className="dashboard-alert" role="alert">
//                     <i className="bi bi-exclamation-circle"></i>

//                     <span>{error}</span>

//                     <button
//                         type="button"
//                         onClick={() => setError("")}
//                         aria-label="Close error message"
//                     >
//                         <i className="bi bi-x-lg"></i>
//                     </button>
//                 </div>
//             )}

//             <div className="create-project-card">
//                 <div className="create-project-heading">
//                     <div className="create-project-icon">
//                         <i className="bi bi-plus-lg"></i>
//                     </div>

//                     <div>
//                         <h2>Create New Project</h2>

//                         <p>
//                             Add a project to start reviewing your source code.
//                         </p>
//                     </div>
//                 </div>

//                 <form
//                     onSubmit={createProject}
//                     className="create-project-form"
//                 >
//                     <div className="project-name-field">
//                         <label htmlFor="projectName">
//                             Project Name
//                         </label>

//                         <div className="dashboard-input-wrapper">
//                             <i className="bi bi-folder"></i>

//                             <input
//                                 id="projectName"
//                                 name="projectName"
//                                 type="text"
//                                 className="dashboard-input"
//                                 placeholder="Enter project name"
//                                 value={name}
//                                 onChange={(event) =>
//                                     setName(event.target.value)
//                                 }
//                                 disabled={creating}
//                                 autoComplete="off"
//                                 required
//                             />
//                         </div>
//                     </div>

//                     <div className="project-description-field">
//                         <label htmlFor="projectDescription">
//                             Description
//                         </label>

//                         <div className="dashboard-input-wrapper">
//                             <i className="bi bi-text-left"></i>

//                             <input
//                                 id="projectDescription"
//                                 name="projectDescription"
//                                 type="text"
//                                 className="dashboard-input"
//                                 placeholder="Short project description"
//                                 value={description}
//                                 onChange={(event) =>
//                                     setDescription(event.target.value)
//                                 }
//                                 disabled={creating}
//                                 autoComplete="off"
//                                    />
//                         </div>
//                     </div>

//                     <button
//                         type="submit"
//                         className="create-project-button"
//                         disabled={creating || !name.trim()}
//                     >
//                         {creating ? (
//                             <>
//                                 <span
//                                     className="spinner-border spinner-border-sm"
//                                     role="status"
//                                     aria-hidden="true"
//                                 ></span>
//                                 <span>Creating</span>
//                             </>
//                         ) : (
//                             <>
//                                 <i className="bi bi-plus-lg"></i>
//                                 <span>Create Project</span>
//                             </>
//                         )}
//                     </button>
//                 </form>
//             </div>

//             <div className="create-project-card mt-4">
//                 <div className="create-project-heading">
//                     <div className="create-project-icon">
//                         <i className="bi bi-lightbulb"></i>
//                     </div>
//                     <div>
//                         <h2>Generate app from prompt</h2>
//                         <p>Build a new frontend/backend app from a requirement prompt and deploy-ready templates.</p>
//                     </div>
//                 </div>
//                 <Link className="create-project-button" to="/generate">
//                     <i className="bi bi-arrow-right-circle"></i>
//                     <span>Open app generator</span>
//                 </Link>
//             </div>

//             {generatedApps.length > 0 && (
//                 <div className="create-project-card mt-4">
//                     <div className="create-project-heading">
//                         <div className="create-project-icon">
//                             <i className="bi bi-layers"></i>
//                         </div>
//                         <div>
//                             <h2>Saved generated apps</h2>
//                             <p>Projects created from prompts and stored in the database are listed here.</p>
//                         </div>
//                     </div>

//                     <div className="project-list-body">
//                         {generatedApps.map((project) => (
//                             <div className="project-list-row" key={project.id}>
//                                 <div className="project-column">
//                                     <div className="project-list-icon">
//                                         <i className="bi bi-code-square"></i>
//                                     </div>
//                                     <div className="project-info">
//                                         <div className="project-name">{project.projectName}</div>
//                                         <span className="project-id">{project.slug}</span>
//                                     </div>
//                                 </div>
//                                 <div className="description-column">
//                                     <span>{project.description || "Generated from prompt"}</span>
//                                 </div>
//                                 <div className="created-column">
//                                     <i className="bi bi-git"></i>
//                                     <span>{project.gitUrl ? "Git linked" : "Local only"}</span>
//                                 </div>
//                                 <div className="action-column">
//                                     <Link to="/generate" className="open-project-button">
//                                         <span>Open</span>
//                                         <i className="bi bi-arrow-right"></i>
//                                     </Link>
//                                 </div>
//                             </div>
//                         ))}
//                     </div>
//                 </div>
//             )}

//             <div className="projects-section">
//                 <div className="projects-section-header">
//                     <div>
//                         <h2>Your Projects</h2>

//                         <p>
//                             Select a project to upload or review source code.
//                         </p>
//                     </div>

//                     <div className="projects-count">
//                         {items.length}{" "}
//                         {items.length === 1 ? "Project" : "Projects"}
//                     </div>
//                 </div>

//                 {loading ? (
//                     <div className="projects-loading">
//                         <div
//                             className="spinner-border text-primary"
//                             role="status"
//                         >
//                             <span className="visually-hidden">
//                                 Loading projects
//                             </span>
//                         </div>

//                         <span>Loading projects...</span>
//                     </div>
//                 ) : items.length === 0 ? (
//                     <div className="projects-empty-state">
//                         <div className="empty-state-icon">
//                             <i className="bi bi-folder2-open"></i>
//                         </div>

//                         <h3>No projects yet</h3>

//                         <p>
//                             Create your first project above to start analyzing
//                             your source code.
//                         </p>
//                     </div>
//                 ) : (
//                     <div className="project-list">
//                         <div className="project-list-header">
//                             <div>Project</div>
//                             <div>Description</div>
//                             <div>Created</div>
//                             <div>Action</div>
//                         </div>

//                         <div className="project-list-body">
//                             {items.map((project) => (
//                                 <div
//                                     className="project-list-row"
//                                     key={project.id}
//                                 >
//                                     <div className="project-column">
//                                         <div className="project-list-icon">
//                                             <i className="bi bi-folder2"></i>
//                                         </div>

//                                         <div className="project-info">
//                                             <Link
//                                                 to={`/projects/${project.id}`}
//                                                 className="project-name"
//                                             >
//                                                 {project.name}
//                                             </Link>

//                                             <span className="project-id">
//                                                 Project #{project.id}
//                                             </span>
//                                         </div>
//                                     </div>

//                                     <div className="description-column">
//                                         {project.description ? (
//                                             <span>
//                                                 {project.description}
//                                             </span>
//                                         ) : (
//                                             <span className="no-description">
//                                                 No description
//                                             </span>
//                                         )}
//                                     </div>

//                                     <div className="created-column">
//                                         <i className="bi bi-calendar3"></i>

//                                         <span>
//                                             {formatCreatedDate(
//                                                 project.createdAt
//                                             )}
//                                         </span>
//                                     </div>

//                                     <div className="action-column">
//                                         <Link
//                                             to={`/projects/${project.id}`}
//                                             className="open-project-button"
//                                             aria-label={`Open ${project.name}`}
//                                         >
//                                             <span>Open</span>
//                                             <i className="bi bi-arrow-right"></i>
//                                         </Link>
//                                     </div>
//                                 </div>
//                             ))}
//                         </div>
//                     </div>
//                 )}
//             </div>
//         </div>
//     );
// };

// export default DashboardPage;
import React, {
    useCallback,
    useEffect,
    useState,
} from "react";
import { Link, useNavigate } from "react-router-dom";
import api from "../api";
import "./DashboardPage.css";

const createEmptyGeneratorForm = () => ({
    projectName: "",
    requirements: "",
    frontendTechnology: "react",
    backendTechnology: "golang",
    databaseTechnology: "postgresql",
    gitUrl: "",
    gitBranch: "main",
});

const DashboardPage = () => {
    const navigate = useNavigate();

    const [items, setItems] = useState([]);
    const [generatedApps, setGeneratedApps] = useState([]);

    const [name, setName] = useState("");
    const [description, setDescription] = useState("");

    const [loading, setLoading] = useState(true);
    const [creating, setCreating] = useState(false);
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");

    const [showGeneratorModal, setShowGeneratorModal] =
        useState(false);

    const [generatorForm, setGeneratorForm] = useState(
        createEmptyGeneratorForm()
    );

    const [generating, setGenerating] = useState(false);
    const [generatorError, setGeneratorError] = useState("");

    const loadProjects = useCallback(async () => {
        try {
            setLoading(true);
            setError("");

            const results = await Promise.all([
                api.get("/api/projects"),
                api.get("/api/projects/generated"),
            ]);

            const normalProjects = Array.isArray(
                results[0].data
            )
                ? results[0].data
                : [];

            const promptProjects = Array.isArray(
                results[1].data
            )
                ? results[1].data
                : [];

            setItems(normalProjects);
            setGeneratedApps(promptProjects);
        } catch (requestError) {
            console.error(
                "Failed to load projects:",
                requestError
            );

            setError(
                "Unable to load projects. Please try again."
            );
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadProjects();
    }, [loadProjects]);

    useEffect(() => {
        if (!showGeneratorModal) {
            return undefined;
        }

        const handleEscapeKey = (event) => {
            if (
                event.key === "Escape" &&
                !generating
            ) {
                setShowGeneratorModal(false);
                setGeneratorError("");
            }
        };

        document.addEventListener(
            "keydown",
            handleEscapeKey
        );

        document.body.classList.add(
            "dashboard-modal-open"
        );

        return () => {
            document.removeEventListener(
                "keydown",
                handleEscapeKey
            );

            document.body.classList.remove(
                "dashboard-modal-open"
            );
        };
    }, [showGeneratorModal, generating]);

    const createProject = async (event) => {
        event.preventDefault();

        const trimmedName = name.trim();
        const trimmedDescription =
            description.trim();

        if (!trimmedName) {
            setError("Project name is required.");
            return;
        }

        try {
            setCreating(true);
            setError("");
            setSuccess("");

            await api.post("/api/projects", {
                name: trimmedName,
                description: trimmedDescription,
            });

            setName("");
            setDescription("");
            setSuccess(
                "Project created successfully."
            );

            await loadProjects();
        } catch (requestError) {
            console.error(
                "Failed to create project:",
                requestError
            );

            setError(
                "Unable to create project. Please try again."
            );
        } finally {
            setCreating(false);
        }
    };

    const openGeneratorModal = () => {
        setGeneratorForm(
            createEmptyGeneratorForm()
        );

        setGeneratorError("");
        setError("");
        setSuccess("");
        setShowGeneratorModal(true);
    };

    const closeGeneratorModal = () => {
        if (generating) {
            return;
        }

        setShowGeneratorModal(false);
        setGeneratorError("");

        setGeneratorForm(
            createEmptyGeneratorForm()
        );
    };

    const handleGeneratorInputChange = (event) => {
        const fieldName = event.target.name;
        const fieldValue = event.target.value;

        setGeneratorForm((currentForm) => ({
            ...currentForm,
            [fieldName]: fieldValue,
        }));
    };

    const generateApplication = async (event) => {
        event.preventDefault();

        const projectName =
            generatorForm.projectName.trim();

        const requirements =
            generatorForm.requirements.trim();

        if (!projectName) {
            setGeneratorError(
                "Please enter a project name."
            );
            return;
        }

        if (!requirements) {
            setGeneratorError(
                "Please enter the application requirements."
            );
            return;
        }

        try {
            setGenerating(true);
            setGeneratorError("");
            setError("");
            setSuccess("");

            await api.post("/api/projects/generate", {
                projectName,
                requirements,
                frontendTechnology:
                    generatorForm.frontendTechnology,
                backendTechnology:
                    generatorForm.backendTechnology,
                databaseTechnology:
                    generatorForm.databaseTechnology,
                gitUrl: generatorForm.gitUrl.trim(),
                gitBranch:
                    generatorForm.gitBranch.trim() ||
                    "main",
            });

            setShowGeneratorModal(false);

            setGeneratorForm(
                createEmptyGeneratorForm()
            );

            setSuccess(
                "Application generated and saved successfully."
            );

            await loadProjects();
        } catch (requestError) {
            console.error(
                "Application generation failed:",
                requestError
            );

            const backendMessage =
                requestError?.response?.data?.error ||
                requestError?.response?.data?.message;

            setGeneratorError(
                backendMessage ||
                    "Unable to generate the application. Check the backend logs and database connection."
            );
        } finally {
            setGenerating(false);
        }
    };

    // const openGeneratedProject = (
    //     projectId
    // ) => {
    //     if (!projectId) {
    //         return;
    //     }

    //     navigate(
    //         `/generate/${encodeURIComponent(
    //             projectId
    //         )}`
    //     );
    // };


    const openGeneratedProject = (projectId) => {
    if (!projectId) {
        return;
    }

    navigate(`/generate/${encodeURIComponent(projectId)}`);
};

    const formatCreatedDate = (createdAt) => {
        if (!createdAt) {
            return "Unknown";
        }

        const parsedDate = new Date(createdAt);

        if (Number.isNaN(parsedDate.getTime())) {
            return "Unknown";
        }

        return parsedDate.toLocaleDateString(
            undefined,
            {
                day: "2-digit",
                month: "short",
                year: "numeric",
            }
        );
    };

    const getTechnologyStack = (project) => {
        const frontend =
            project.frontendTechnology ||
            "Frontend";

        const backend =
            project.backendTechnology ||
            "Backend";

        const database =
            project.databaseTechnology ||
            "No database";

        return `${frontend} + ${backend} + ${database}`;
    };

    return (
        <div className="dashboard-page">
            <div className="dashboard-header">
                <div>
                    <div className="dashboard-label">
                        <i className="bi bi-grid"></i>
                        <span>Dashboard</span>
                    </div>

                    <h1>Projects</h1>

                    <p>
                        Create and manage projects,
                        upload source code, generate
                        applications from requirements,
                        or start a review from a Git
                        repository.
                    </p>
                </div>

                <div className="dashboard-summary">
                    <div className="summary-icon">
                        <i className="bi bi-folder2-open"></i>
                    </div>

                    <div>
                        <span>Total Projects</span>

                        <strong>
                            {items.length +
                                generatedApps.length}
                        </strong>
                    </div>
                </div>
            </div>

            {error && (
                <div
                    className="dashboard-alert dashboard-alert-error"
                    role="alert"
                >
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
                <div
                    className="dashboard-alert dashboard-alert-success"
                    role="status"
                >
                    <i className="bi bi-check-circle"></i>

                    <span>{success}</span>

                    <button
                        type="button"
                        onClick={() => setSuccess("")}
                        aria-label="Close success message"
                    >
                        <i className="bi bi-x-lg"></i>
                    </button>
                </div>
            )}

            <section className="dashboard-card">
                <div className="dashboard-card-heading">
                    <div className="dashboard-card-icon">
                        <i className="bi bi-plus-lg"></i>
                    </div>

                    <div>
                        <h2>Create New Project</h2>

                        <p>
                            Add a project to start
                            reviewing your source code.
                        </p>
                    </div>
                </div>

                <form
                    className="create-project-form"
                    onSubmit={createProject}
                >
                    <div className="dashboard-field">
                        <label htmlFor="projectName">
                            Project Name
                        </label>

                        <div className="dashboard-input-wrapper">
                            <i className="bi bi-folder"></i>

                            <input
                                id="projectName"
                                name="projectName"
                                type="text"
                                value={name}
                                onChange={(event) =>
                                    setName(
                                        event.target.value
                                    )
                                }
                                placeholder="Enter project name"
                                disabled={creating}
                                autoComplete="off"
                                required
                            />
                        </div>
                    </div>

                    <div className="dashboard-field">
                        <label htmlFor="projectDescription">
                            Description
                        </label>

                        <div className="dashboard-input-wrapper">
                            <i className="bi bi-text-left"></i>

                            <input
                                id="projectDescription"
                                name="projectDescription"
                                type="text"
                                value={description}
                                onChange={(event) =>
                                    setDescription(
                                        event.target.value
                                    )
                                }
                                placeholder="Short project description"
                                disabled={creating}
                                autoComplete="off"
                            />
                        </div>
                    </div>

                    <button
                        type="submit"
                        className="dashboard-primary-button"
                        disabled={
                            creating || !name.trim()
                        }
                    >
                        {creating ? (
                            <>
                                <span
                                    className="spinner-border spinner-border-sm"
                                    aria-hidden="true"
                                ></span>

                                <span>Creating</span>
                            </>
                        ) : (
                            <>
                                <i className="bi bi-plus-lg"></i>
                                <span>Create Project</span>
                            </>
                        )}
                    </button>
                </form>
            </section>

            <section className="generator-launch-card">
                <div className="generator-launch-content">
                    <div className="generator-launch-icon">
                        <i className="bi bi-stars"></i>
                    </div>

                    <div>
                        <h2>
                            Generate App From Prompt
                        </h2>

                        <p>
                            Build a new frontend and
                            backend application from a
                            requirement prompt using
                            deploy-ready project templates.
                        </p>
                    </div>
                </div>

                <button
                    type="button"
                    className="dashboard-primary-button"
                    onClick={openGeneratorModal}
                >
                    <i className="bi bi-magic"></i>

                    <span>
                        Open App Generator
                    </span>
                </button>
            </section>

            <section className="projects-section generated-projects-section">
                <div className="projects-section-header">
                    <div>
                        <h2>
                            Saved Generated Apps
                        </h2>

                        <p>
                            Applications created from
                            prompts and saved in the
                            database.
                        </p>
                    </div>

                    <div className="projects-count">
                        {generatedApps.length}{" "}
                        {generatedApps.length === 1
                            ? "Application"
                            : "Applications"}
                    </div>
                </div>

                {loading ? (
                    <div className="projects-loading">
                        <div
                            className="spinner-border"
                            role="status"
                        >
                            <span className="visually-hidden">
                                Loading generated
                                applications
                            </span>
                        </div>

                        <span>
                            Loading generated
                            applications...
                        </span>
                    </div>
                ) : generatedApps.length === 0 ? (
                    <div className="projects-empty-state">
                        <div className="empty-state-icon">
                            <i className="bi bi-stars"></i>
                        </div>

                        <h3>
                            No generated applications
                        </h3>

                        <p>
                            Use Open App Generator to
                            create an application from
                            requirements.
                        </p>
                    </div>
                ) : (
                    <div className="project-table">
                        <div className="project-table-header">
                            <div>Application</div>
                            <div>
                                Technology Stack
                            </div>
                            <div>Git Status</div>
                            <div>Action</div>
                        </div>

                        <div className="project-table-body">
                            {generatedApps.map(
                                (project) => (
                                    <div
                                        className="project-table-row"
                                        key={project.id}
                                    >
                                        <div className="project-column">
                                            <div className="project-list-icon">
                                                <i className="bi bi-code-square"></i>
                                            </div>

                                            <div className="project-info">
                                                <button
                                                    type="button"
                                                    className="project-name-button"
                                                    onClick={() =>
                                                        openGeneratedProject(
                                                            project.id
                                                        )
                                                    }
                                                >
                                                    {project.projectName ||
                                                        "Generated Application"}
                                                </button>

                                                <span className="project-id">
                                                    {
                                                        project.id
                                                    }
                                                </span>
                                            </div>
                                        </div>

                                        <div className="description-column">
                                            {getTechnologyStack(
                                                project
                                            )}
                                        </div>

                                        <div className="created-column">
                                            <i className="bi bi-git"></i>

                                            <span>
                                                {project.gitUrl
                                                    ? "Git linked"
                                                    : "Local only"}
                                            </span>
                                        </div>

                                        <div className="action-column">
                                            <button
                                                type="button"
                                                className="open-project-button"
                                                onClick={() =>
                                                    openGeneratedProject(
                                                        project.id
                                                    )
                                                }
                                            >
                                                <span>
                                                    Open
                                                </span>

                                                <i className="bi bi-arrow-right"></i>
                                            </button>
                                        </div>
                                    </div>
                                )
                            )}
                        </div>
                    </div>
                )}
            </section>

            <section className="projects-section">
                <div className="projects-section-header">
                    <div>
                        <h2>Your Projects</h2>

                        <p>
                            Select a project to upload or
                            review source code.
                        </p>
                    </div>

                    <div className="projects-count">
                        {items.length}{" "}
                        {items.length === 1
                            ? "Project"
                            : "Projects"}
                    </div>
                </div>

                {loading ? (
                    <div className="projects-loading">
                        <div
                            className="spinner-border"
                            role="status"
                        >
                            <span className="visually-hidden">
                                Loading projects
                            </span>
                        </div>

                        <span>
                            Loading projects...
                        </span>
                    </div>
                ) : items.length === 0 ? (
                    <div className="projects-empty-state">
                        <div className="empty-state-icon">
                            <i className="bi bi-folder2-open"></i>
                        </div>

                        <h3>No projects yet</h3>

                        <p>
                            Create your first project to
                            start analyzing source code.
                        </p>
                    </div>
                ) : (
                    <div className="project-table">
                        <div className="project-table-header">
                            <div>Project</div>
                            <div>Description</div>
                            <div>Created</div>
                            <div>Action</div>
                        </div>

                        <div className="project-table-body">
                            {items.map((project) => (
                                <div
                                    className="project-table-row"
                                    key={project.id}
                                >
                                    <div className="project-column">
                                        <div className="project-list-icon">
                                            <i className="bi bi-folder2"></i>
                                        </div>

                                        <div className="project-info">
                                            <Link
                                                to={`/projects/${project.id}`}
                                                className="project-name-link"
                                            >
                                                {project.name}
                                            </Link>

                                            <span className="project-id">
                                                Project #
                                                {project.id}
                                            </span>
                                        </div>
                                    </div>

                                    <div className="description-column">
                                        {project.description ||
                                            "No description"}
                                    </div>

                                    <div className="created-column">
                                        <i className="bi bi-calendar3"></i>

                                        <span>
                                            {formatCreatedDate(
                                                project.createdAt
                                            )}
                                        </span>
                                    </div>

                                    <div className="action-column">
                                        <Link
                                            to={`/projects/${project.id}`}
                                            className="open-project-button"
                                        >
                                            <span>
                                                Open
                                            </span>

                                            <i className="bi bi-arrow-right"></i>
                                        </Link>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </section>

            {showGeneratorModal && (
                <div className="generator-modal-overlay">
                    <div
                        className="generator-modal"
                        role="dialog"
                        aria-modal="true"
                        aria-labelledby="generatorModalTitle"
                    >
                        <div className="generator-modal-header">
                            <div className="generator-modal-heading">
                                <div className="generator-modal-icon">
                                    <i className="bi bi-stars"></i>
                                </div>

                                <div>
                                    <h2 id="generatorModalTitle">
                                        Generate Application
                                        From Requirements
                                    </h2>

                                    <p>
                                        Select the
                                        technology stack and
                                        describe the
                                        application.
                                    </p>
                                </div>
                            </div>

                            <button
                                type="button"
                                className="generator-close-button"
                                onClick={
                                    closeGeneratorModal
                                }
                                disabled={generating}
                                aria-label="Close generator"
                            >
                                <i className="bi bi-x-lg"></i>
                            </button>
                        </div>

                        <form
                            className="generator-modal-body"
                            onSubmit={
                                generateApplication
                            }
                        >
                            {generatorError && (
                                <div
                                    className="generator-form-error"
                                    role="alert"
                                >
                                    <i className="bi bi-exclamation-circle"></i>

                                    <span>
                                        {generatorError}
                                    </span>
                                </div>
                            )}

                            <div className="generator-field">
                                <label htmlFor="generatedProjectName">
                                    Project Name
                                </label>

                                <input
                                    id="generatedProjectName"
                                    name="projectName"
                                    type="text"
                                    value={
                                        generatorForm.projectName
                                    }
                                    onChange={
                                        handleGeneratorInputChange
                                    }
                                    placeholder="For example: kids-shop"
                                    disabled={generating}
                                    autoComplete="off"
                                    autoFocus
                                    required
                                />
                            </div>

                            <div className="generator-stack-grid">
                                <div className="generator-field">
                                    <label htmlFor="frontendTechnology">
                                        Frontend
                                    </label>

                                    <select
                                        id="frontendTechnology"
                                        name="frontendTechnology"
                                        value={
                                            generatorForm.frontendTechnology
                                        }
                                        onChange={
                                            handleGeneratorInputChange
                                        }
                                        disabled={generating}
                                    >
                                        <option value="react">
                                            React
                                        </option>

                                        <option value="angular">
                                            Angular
                                        </option>
                                    </select>
                                </div>

                                <div className="generator-field">
                                    <label htmlFor="backendTechnology">
                                        Backend
                                    </label>

                                    <select
                                        id="backendTechnology"
                                        name="backendTechnology"
                                        value={
                                            generatorForm.backendTechnology
                                        }
                                        onChange={
                                            handleGeneratorInputChange
                                        }
                                        disabled={generating}
                                    >
                                        <option value="golang">
                                            Go
                                        </option>

                                        <option value="java">
                                            Java Spring Boot
                                        </option>
                                    </select>
                                </div>

                   <div className="generator-field">
                                    <label htmlFor="databaseTechnology">
                                        Database
                                    </label>

                                    <select
                                        id="databaseTechnology"
                                        name="databaseTechnology"
                                        value={
                                            generatorForm.databaseTechnology
                                        }
                                        onChange={
                                            handleGeneratorInputChange
                                        }
                                        disabled={generating}
                                    >
                                        <option value="postgresql">
                                            PostgreSQL
                                        </option>

                                        <option value="mysql">
                                            MySQL
                                        </option>

                                        <option value="mongodb">
                                            MongoDB
                                        </option>

                                        <option value="none">
                                            No Database
                                        </option>
                                    </select>
                                </div>
                            </div>

                            <div className="generator-field">
                                <label htmlFor="requirements">
                                    Application
                                    Requirements
                                </label>

                                <textarea
                                    id="requirements"
                                    name="requirements"
                                    rows={9}
                                    value={
                                        generatorForm.requirements
                                    }
                                    onChange={
                                        handleGeneratorInputChange
                                    }
                                    placeholder="Describe pages, authentication, product features, workflows, API requirements, database entities and deployment requirements."
                                    disabled={generating}
                                    required
                                />
                            </div>

                            <div className="generator-git-grid">
                                <div className="generator-field">
                                    <label htmlFor="gitUrl">
                                        Git Remote URL
                                    </label>

                                    <input
                                        id="gitUrl"
                                        name="gitUrl"
                                        type="text"
                                        value={
                                            generatorForm.gitUrl
                                        }
                                        onChange={
                                            handleGeneratorInputChange
                                        }
                                        placeholder="Optional Git repository URL"
                                        disabled={generating}
                                        autoComplete="off"
                                    />
                                </div>

                                <div className="generator-field">
                                    <label htmlFor="gitBranch">
                                        Git Branch
                                    </label>

                                    <input
                                        id="gitBranch"
                                        name="gitBranch"
                                        type="text"
                                        value={
                                            generatorForm.gitBranch
                                        }
                                        onChange={
                                            handleGeneratorInputChange
                                        }
                                        placeholder="main"
                                        disabled={generating}
                                        autoComplete="off"
                                    />
                                </div>
                            </div>

                            <div className="generator-modal-actions">
                                <button
                                    type="button"
                                    className="generator-cancel-button"
                                    onClick={
                                        closeGeneratorModal
                                    }
                                    disabled={generating}
                                >
                                    Cancel
                                </button>

                                <button
                                    type="submit"
                                    className="generator-submit-button"
                                    disabled={
                                        generating ||
                                        !generatorForm.projectName.trim() ||
                                        !generatorForm.requirements.trim()
                                    }
                                >
                                    {generating ? (
                                        <>
                                            <span
                                                className="spinner-border spinner-border-sm"
                                                aria-hidden="true"
                                            ></span>

                                            <span>
                                                Generating
                                            </span>
                                        </>
                                    ) : (
                                        <>
                                            <i className="bi bi-magic"></i>

                                            <span>
                                                Generate
                                                Application
                                            </span>
                                        </>
                                    )}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default DashboardPage;