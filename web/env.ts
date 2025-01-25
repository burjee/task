const dev = process.env.NODE_ENV === "development";

export const HTTPURL = dev ? "http://localhost:8000/" : "/";