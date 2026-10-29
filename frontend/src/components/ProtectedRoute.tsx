import React from 'react';import{Navigate}from'react-router-dom';import{useAuth}from'../auth/AuthContext';
const ProtectedRoute:React.FC<React.PropsWithChildren>=({children})=>{const{user,loading}=useAuth();if(loading)return <div className="p-5 text-center">Loading...</div>;return user?<>{children}</>:<Navigate to="/login" replace/>};export default ProtectedRoute;
