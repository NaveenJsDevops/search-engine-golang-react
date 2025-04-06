const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:9000';

export const API_ENDPOINTS = {
    FETCH_LOGS: `${BASE_URL}/v1/list/log/entries`,
    UPLOAD_PARQUET: `${BASE_URL}/v1/upload/perquet?imageFor=new`,
    FETCH_TOTAL_LOGS: `${BASE_URL}/v1/fetch/all/records`,
};
