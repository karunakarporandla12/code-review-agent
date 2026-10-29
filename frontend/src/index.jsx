// import React from "react";
// import ReactDOM from "react-dom/client";

// import "bootstrap/dist/css/bootstrap.min.css";
// import "bootstrap-icons/font/bootstrap-icons.css";

// import "./styles.css";
// import App from "./App";

// const rootElement = document.getElementById("root");

// if (!rootElement) {
//     throw new Error("Root element was not found.");
// }

// const root = ReactDOM.createRoot(rootElement);

// root.render(
//     <React.StrictMode>
//         <App />
//     </React.StrictMode>
// );

import React from "react";
import ReactDOM from "react-dom/client";

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-icons/font/bootstrap-icons.css";

import "./styles.css";
import "./styles/Application.css";
import App from "./App";

const rootElement = document.getElementById("root");

if (!rootElement) {
    throw new Error("Root element not found");
}

const root = ReactDOM.createRoot(rootElement);

root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);