export interface User { id: string; name: string; email: string; role: string; }
export interface Project { id: string; ownerId: string; name: string; description: string; createdAt: string; }
export interface Review { id:string; projectId:string; sourceType:string; sourceName:string; branchName?:string; status:string; summary:string; filesScanned:number; criticalCount:number; highCount:number; mediumCount:number; lowCount:number; errorMessage:string; createdAt:string; }
export interface Finding { id:string; severity:string; category:string; ruleId:string; filePath:string; lineNumber:number; title:string; description:string; recommendation:string; codeExcerpt:string; source:string; status:string; }
