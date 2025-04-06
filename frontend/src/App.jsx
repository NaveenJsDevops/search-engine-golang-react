import React, { useState, useEffect } from 'react';
import axios from "axios";
import { toast } from "react-toastify";
import ResultsTable from "./components/ResultsTable";
import { API_ENDPOINTS } from '../src/components/apiUrls';

function App() {
    const [results, setResults] = useState([]);
    const [stats, setStats] = useState({ time: 0, count: 0 });
    const [pagination, setPagination] = useState({ limit: 10, offset: 0, currentPage: 1, total: 0 });
    const [query, setQuery] = useState("");
    const [selectedFile, setSelectedFile] = useState(null);
    const [type, setType] = useState("TYPE_1"); // TYPE_1 = paginated, TYPE_2 = full list
    const [loading, setLoading] = useState(false);

    const fetchData = async () => {
        if (!query.trim()) return;

        try {
            const url = type === "TYPE_1" ? API_ENDPOINTS.FETCH_LOGS : API_ENDPOINTS.FETCH_TOTAL_LOGS;
            const params = type === "TYPE_1"
                ? { limit: pagination.limit, offset: pagination.offset, search: query }
                : { search: query };

            const response = await axios.get(url, { params });
            const logs = response.data?.log || [];

            setResults(logs);
            setStats({
                count: logs.length,
                time: response.data?.duration || 0
            });

            if (type === "TYPE_1") {
                setPagination(prev => ({ ...prev, total: response.data?.total || 0 }));
            }
        } catch (error) {
            console.error("Error fetching data:", error);
            toast.error('Failed to fetch data');
        }
    };


    const handleSearch = () => {
        if (!query.trim()) {
            toast.warning("Please enter a search keyword.");
            return;
        }
        setPagination({ ...pagination, offset: 0 });
        fetchData();
    };

    const handleFileChange = (event) => {
        const file = event.target.files[0];
        if (file) setSelectedFile(file);
    };

    const handleUpload = async () => {
        if (!selectedFile) {
            toast.error("Please select a file first.");
            return;
        }
        const formData = new FormData();
        formData.append("image", selectedFile);
        try {
            await axios.post(API_ENDPOINTS.UPLOAD_PARQUET, formData);
            toast.success("File uploaded successfully");
            setSelectedFile(null);
        } catch (error) {
            console.error("Upload failed:", error);
            toast.error("Upload failed");
        }
    };

    useEffect(() => {
        if (query.trim()) {
            fetchData();
        }
    }, [pagination.limit, pagination.offset, type]);

    return (
        <div style={{ fontFamily: "sans-serif", padding: 15, background: "#f9fafb" }}>
            {loading && (
                <div style={{
                    position: "fixed",
                    top: 0,
                    left: 0,
                    width: "100vw",
                    height: "100vh",
                    backgroundColor: "rgba(255, 255, 255, 0.8)",
                    zIndex: 9999,
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    fontSize: "20px",
                    fontWeight: "bold",
                    color: "#1e3a8a"
                }}>
                    ⏳ Loading records, please wait...
                </div>
            )}
            <h1 style={{ textAlign: "center", color: "#1e3a8a", fontSize: 30, marginBottom: 20 }}>
                🔍 Apica In-Memory Search Engine
            </h1>

            {/* Type Toggle */}
            <div style={{ position: 'absolute', top: 15, left: 20, fontWeight: "bold" }}>
                <label>
                    <input
                        type="radio"
                        name="type"
                        value="TYPE_1"
                        checked={type === "TYPE_1"}
                        onChange={() => {
                            setResults([]);
                            setPagination({ ...pagination, offset: 0, currentPage: 1 });
                            setType("TYPE_1");
                        }}

                    /> Type 1
                </label>
                {"  "}
                <label style={{ marginLeft: "10px" }}>
                    <input
                        type="radio"
                        name="type"
                        value="TYPE_2"
                        checked={type === "TYPE_2"}
                        onChange={() => {
                            setResults([]);
                            setPagination({ ...pagination, offset: 0, currentPage: 1 });
                            setType("TYPE_2");
                        }}

                    /> Type 2
                </label>
            </div>

            {/* Search Bar */}
            <div style={{ display: "flex", justifyContent: "center", alignItems: "center" }}>
                <input
                    type="text"
                    placeholder="Search keywords..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
                    style={{ padding: '8px', margin: '0 5px', border: '1px solid #ccc', borderRadius: '4px', width: '300px' }}
                />
                <button onClick={handleSearch} style={{
                    backgroundColor: '#4CAF50', color: 'white', padding: '8px 15px', margin: '0 5px',
                    border: 'none', borderRadius: '4px', cursor: 'pointer', fontWeight: "bold"
                }}>
                    Search
                </button>
            </div>

            {/* Upload Section */}
            <div style={{
                display: "flex",
                justifyContent: "flex-end",
                marginTop: "-80px",
                paddingRight: "20px"
            }}>
                <div style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "10px",
                    flexWrap: "wrap"
                }}>
                    <label style={{ fontWeight: "bold", fontSize: "14px", whiteSpace: "nowrap" }}>
                        Upload Parquet File:
                    </label>
                    <div
                        onClick={() => document.getElementById('hidden-file').click()}
                        style={{
                            padding: '4px 8px',
                            border: '1px solid #ccc',
                            borderRadius: '6px',
                            backgroundColor: '#fff',
                            cursor: 'pointer',
                            fontStyle: selectedFile ? 'normal' : 'italic',
                            color: selectedFile ? '#000' : '#666',
                            display: 'flex',
                            alignItems: 'center',
                            gap: '8px',
                            boxShadow: '1px 1px 5px rgba(0,0,0,0.1)',
                            whiteSpace: "nowrap",
                            fontSize: "12px",
                        }}
                    >
                        📁 {selectedFile ? selectedFile.name : 'Choose File'}
                    </div>
                    <input id="hidden-file" type="file" onChange={handleFileChange} style={{ display: 'none' }} />
                    <button onClick={handleUpload} style={{
                        backgroundColor: '#4CAF50',
                        color: 'white',
                        padding: '8px 15px',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        fontWeight: "bold"
                    }}>
                        Upload
                    </button>
                </div>
            </div>

            {/* Stats */}
            <div style={{
                position: 'fixed',
                bottom: '2px',
                right: '50px',
                textAlign: 'right',
                backgroundColor: '#f9fafb',
                padding: '10px 15px',
                borderRadius: '6px',
                zIndex: 1001,
                display: 'flex',
                gap: '20px',
                fontWeight: 'bold',
                fontSize: '15px',
            }}>
                <span>Total Found = {stats.count}</span>|<span>Duration = {stats.time}</span>
            </div>

            {/* Table */}
            <ResultsTable
                data={results}
                pagination={pagination}
                setPagination={setPagination}
                fetchData={fetchData}
                hidePagination={type === "TYPE_2"}
                query={query}
            />
        </div>
    );
}

export default App;
