// import React from "react";
// import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
// import { AuthProvider } from "./auth/AuthContext";
// import Layout from "./components/Layout";
// import ProtectedRoute from "./components/ProtectedRoute";
// import LoginPage from "./pages/LoginPage";
// import RegisterPage from "./pages/RegisterPage";
// import ProjectPage from "./pages/ProjectPage";
// import FeatureGeneratorPage from "./pages/FeatureGeneratorPage";
// import ReviewPage from "./pages/ReviewPage";
// import HomePage from "./pages/HomePage";
// import DashboardPage from "./pages/DashboardPage";

// const App = () => (
//   <BrowserRouter>
//     <AuthProvider>
//       <Layout>
//         <Routes>
//           <Route path="/" element={<HomePage />} />
//           <Route path="/login" element={<LoginPage />} />
//           <Route path="/register" element={<RegisterPage />} />
//           <Route
//             path="/dashboard"
//             element={
//               <ProtectedRoute>
//                 <DashboardPage />
//               </ProtectedRoute>
//             }
//           />
//           {/* <Route
//             path="/projects/:projectId"
//             element={
//               <ProtectedRoute>
//                 <ProjectPage />
//               </ProtectedRoute>
//             }
//           /> */}
//           {/* <Route
//             path="/generate"
//             element={
//               <ProtectedRoute>
//                 <FeatureGeneratorPage />
//               </ProtectedRoute>
//             }
//           /> */}
//           {/* <Route
//             path="/projects/:projectId/generate"
//             element={
//               <ProtectedRoute>
//                 <FeatureGeneratorPage />
//               </ProtectedRoute>
//             }
//           /> */}
//           <Route
//             path="/projects/:projectId"
//             element={
//               <ProtectedRoute>
//                 <FeatureGeneratorPage />
//               </ProtectedRoute>
//             }
//           />
//           <Route
//             path="/reviews/:reviewId"
//             element={
//               <ProtectedRoute>
//                 <ReviewPage />
//               </ProtectedRoute>
//             }
//           />
//           <Route path="*" element={<Navigate to="/" />} />
//         </Routes>
//       </Layout>
//     </AuthProvider>
//   </BrowserRouter>
// );

// export default App;


import React from "react";
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";

import { AuthProvider } from "./auth/AuthContext";
import Layout from "./components/Layout";
import ProtectedRoute from "./components/ProtectedRoute";

import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import DashboardPage from "./pages/DashboardPage";
import ProjectPage from "./pages/ProjectPage";
import FeatureGeneratorPage from "./pages/FeatureGeneratorPage";
import ReviewPage from "./pages/ReviewPage";

const App = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Layout>
          <Routes>
            <Route
              path="/"
              element={<HomePage />}
            />

            <Route
              path="/login"
              element={<LoginPage />}
            />

            <Route
              path="/register"
              element={<RegisterPage />}
            />

            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <DashboardPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/projects/:projectId"
              element={
                <ProtectedRoute>
                  <ProjectPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/generate/:projectId"
              element={
                <ProtectedRoute>
                  <FeatureGeneratorPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="/generate"
              element={
                <Navigate
                  to="/dashboard"
                  replace
                />
              }
            />

            <Route
              path="/reviews/:reviewId"
              element={
                <ProtectedRoute>
                  <ReviewPage />
                </ProtectedRoute>
              }
            />

            <Route
              path="*"
              element={
                <Navigate
                  to="/"
                  replace
                />
              }
            />
          </Routes>
        </Layout>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
