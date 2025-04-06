import React, { useState, useEffect } from 'react';

function ResultsTable({ data, pagination, setPagination, fetchData, hidePagination = false, query }) {
    const [showModal, setShowModal] = useState(false);
    const [modalContent, setModalContent] = useState('');
    const [modalTitle, setModalTitle] = useState('');
    const [renderedRows, setRenderedRows] = useState([]);

    const handlePageChange = (newPage) => {
        setPagination(prev => ({ ...prev, currentPage: newPage, offset: (newPage - 1) * prev.limit }));
        fetchData();
    };

    const handleToggleModal = (json, title) => {
        setModalContent(json);
        setModalTitle(title);
        setShowModal(true);
    };

    const highlightText = (text, search) => {
        if (!text || !search) return text;
        const regex = new RegExp(`(${search})`, 'gi');
        return text.toString().replace(regex, `<mark style="background: #ffe58a;">$1</mark>`);
    };

    useEffect(() => {
        if (hidePagination && data.length > 0) {
            const chunkSize = 100;
            let currentIndex = 0;
            const tempList = [];

            const renderChunk = () => {
                const chunk = data.slice(currentIndex, currentIndex + chunkSize);
                tempList.push(...chunk);
                setRenderedRows([...tempList]);

                currentIndex += chunkSize;
                if (currentIndex < data.length) {
                    setTimeout(renderChunk, 0); // Allow UI to breathe
                }
            };

            renderChunk();
        } else {
            setRenderedRows(data);
        }
    }, [data, hidePagination]);

    return (
        <div style={{ display: 'flex', flexDirection: 'column', height: 'calc(100vh - 140px)', overflow: 'hidden' }}>
            <style>
                {`
                  ::-webkit-scrollbar {
                    display: none;
                  }
                `}
            </style>

            {/* Table Section */}
            <div style={{ flexGrow: 1, overflowY: 'auto', overflowX: 'auto', marginTop: '55px' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', tableLayout: 'auto' }}>
                    <thead style={{ position: "sticky", top: 0, backgroundColor: "#64aef5", zIndex: 1 }}>
                    <tr>
                        {[
                            'Msg Id', 'Partition Id', 'Timestamp', 'Hostname', 'Priority', 'Facility',
                            'Facility String', 'Severity', 'Severity String', 'App Name', 'Proc Id',
                            'Message', 'Sender', 'Nano TimeStamp', 'Namespace', 'Message Raw', 'Structured Data'
                        ].map((header, index) => (
                            <th key={index} style={{
                                border: '1px solid #ddd',
                                padding: '10px 10px 5px 10px',
                                backgroundColor: '#64aef5',
                                textAlign: 'left',
                                whiteSpace: 'nowrap',
                                fontSize: '12px',
                                minWidth: '120px',
                            }}>
                                {header}
                            </th>
                        ))}
                    </tr>
                    </thead>
                    <tbody style={{ fontSize: '11px' }}>
                    {renderedRows.map((item, index) => (
                        <tr key={index} style={{
                            borderBottom: '1px solid #ddd',
                            backgroundColor: index % 2 === 0 ? '#ffffff' : '#f7faff'
                        }}>
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.MsgId, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.PartitionId, query) }} />
                            <td style={tdStyle}>
                                {new Date(item.Timestamp).toISOString().replace('T', ' ').substring(0, 19)}
                            </td>
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Hostname, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Priority, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Facility, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.FacilityString, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Severity, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.SeverityString, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.AppName, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.ProcId, query) }} />
                            <td style={{
                                ...tdStyle,
                                maxWidth: '200px',
                                whiteSpace: 'nowrap',
                                overflow: 'hidden',
                                textOverflow: 'ellipsis'
                            }} title={item.Message}
                                dangerouslySetInnerHTML={{ __html: highlightText(item.Message, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Sender, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.NanoTimeStamp, query) }} />
                            <td style={tdStyle} dangerouslySetInnerHTML={{ __html: highlightText(item.Namespace, query) }} />
                            <td style={tdStyle}>
                                <button onClick={() => handleToggleModal(item.MessageRaw, "Message Raw")} style={buttonStyle}>
                                    Show
                                </button>
                            </td>
                            <td style={tdStyle}>
                                <button onClick={() => handleToggleModal(item.StructuredData, "Structured Data")} style={buttonStyle}>
                                    Show
                                </button>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>

            {/* Pagination Section */}
            {!hidePagination && (
                <div style={{
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    width: '100%',
                    backgroundColor: '#f9fafb',
                    borderTop: '1px solid #ccc',
                    padding: '10px 20px',
                    display: 'flex',
                    justifyContent: 'flex-start',
                    alignItems: 'center',
                    zIndex: 999,
                    flexWrap: 'wrap',
                    gap: '16px'
                }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                        {data.length !== 0 && (
                            <button
                                onClick={() => handlePageChange(pagination.currentPage - 1)}
                                disabled={pagination.currentPage === 1}
                                style={navBtnStyle}
                            >
                                ◀ Prev
                            </button>
                        )}
                        {data.length >= pagination.limit && (
                            <button
                                onClick={() => handlePageChange(pagination.currentPage + 1)}
                                style={navBtnStyle}
                            >
                                Next ▶
                            </button>
                        )}
                    </div>
                    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                        <span style={{ fontSize: '14px', fontWeight: 'bold' }}>Showing</span>
                        <select
                            value={pagination.limit}
                            onChange={(e) => {
                                const newLimit = parseInt(e.target.value);
                                setPagination({
                                    ...pagination,
                                    limit: newLimit,
                                    offset: 0,
                                    currentPage: 1
                                });
                                fetchData();
                            }}
                            style={{
                                padding: '6px',
                                borderRadius: '4px',
                                border: '1px solid #ccc',
                                fontSize: '14px',
                                cursor: 'pointer'
                            }}
                        >
                            {[5, 10, 20, 25, 50, 100].map(size => (
                                <option key={size} value={size}>{size}</option>
                            ))}
                        </select>
                        <span style={{ fontSize: '14px', color: '#555' }}>
                            of {pagination.total} records
                        </span>
                    </div>
                </div>
            )}

            {/* JSON Modal */}
            {showModal && (
                <div style={{
                    position: "fixed",
                    top: 0, left: 0,
                    width: "100%", height: "100%",
                    backgroundColor: "rgba(0, 0, 0, 0.5)",
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    zIndex: 1010
                }}>
                    <div style={{
                        background: "#fff",
                        padding: "20px",
                        borderRadius: "8px",
                        maxHeight: "80%",
                        maxWidth: "80%",
                        overflowY: "auto"
                    }}>
                        <h3 style={{ marginBottom: "10px" }}>{modalTitle}</h3>
                        <pre style={{
                            whiteSpace: "pre-wrap",
                            wordWrap: "break-word",
                            backgroundColor: "#f3f4f6",
                            padding: "10px",
                            border: "1px solid #ccc",
                            borderRadius: "6px",
                            fontSize: "13px"
                        }}>
                            {JSON.stringify(JSON.parse(modalContent), null, 2)}
                        </pre>
                        <button onClick={() => setShowModal(false)} style={{
                            marginTop: "10px",
                            padding: "8px 16px",
                            backgroundColor: "#4CAF50",
                            color: "white",
                            border: "none",
                            borderRadius: "4px",
                            cursor: "pointer"
                        }}>Close</button>
                    </div>
                </div>
            )}
        </div>
    );
}

const tdStyle = {
    border: '1px solid #ddd',
    padding: '10px',
    verticalAlign: 'top'
};

const buttonStyle = {
    padding: "4px 10px",
    backgroundColor: "#4CAF50",
    color: "#fff",
    border: "none",
    borderRadius: "4px",
    cursor: "pointer"
};

const navBtnStyle = {
    padding: '6px 10px',
    backgroundColor: '#4CAF50',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
    fontSize: '13px'
};

export default ResultsTable;
