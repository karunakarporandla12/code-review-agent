import React, {createContext,useContext,useEffect,useState} from 'react';
import api from '../api';
import {User} from '../types';
interface AuthValue{user:User|null;loading:boolean;setSession:(token:string,user:User)=>void;logout:()=>void;}
const AuthContext=createContext<AuthValue|undefined>(undefined);
export const AuthProvider:React.FC<React.PropsWithChildren>=({children})=>{const[user,setUser]=useState<User|null>(null);const[loading,setLoading]=useState(true);useEffect(()=>{const token=localStorage.getItem('token');if(!token){setLoading(false);return;}api.get<User>('/api/me').then(r=>setUser(r.data)).catch(()=>localStorage.removeItem('token')).finally(()=>setLoading(false));},[]);const setSession=(token:string,value:User)=>{localStorage.setItem('token',token);setUser(value);};const logout=()=>{localStorage.removeItem('token');setUser(null);};return <AuthContext.Provider value={{user,loading,setSession,logout}}>{children}</AuthContext.Provider>};
export const useAuth=()=>{const value=useContext(AuthContext);if(!value)throw new Error('useAuth must be inside AuthProvider');return value;};
