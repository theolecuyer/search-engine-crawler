"use client";

import { useState } from "react";
import axios from "axios";

export default function Home() {
  const [isLoading, setIsLoading] = useState(false); // State for loading status
  const [searchTerm, setSearchTerm] = useState(""); // To store the searchTerm being searched
  const [results, setResults] = useState([]); // State to hold search results

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const apiKey = process.env.NEXT_PUBLIC_GO_API_KEY || "";

  const handleSearch = async () => {
    setIsLoading(true);
    console.log("handleSearch called");
    console.log("Test API KEY", apiKey);
    try {
      const response = await axios.post(
        apiUrl,
        { search_term: searchTerm },
        {
          headers: {
            "Authorization": apiKey,
          },
        }
      );
      const data = response.data || [];
      setResults(data);
    } catch (error) {
      console.error("Error calling API:", error);
      setResults([]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-black">
      {/* Header with search bar */}
      <div className="flex items-center justify-between p-4 bg-gray-800">
        <h1 className="text-2xl font-bold text-white">SiteSearch</h1>
        <div className="flex gap-2">
          <input
            type="text"
            className="border border-gray-300 rounded-md p-2 w-64 text-black"
            placeholder="Search..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <button
            className="bg-blue-500 text-white rounded px-4 h-full transition duration-300 hover:bg-blue-600"
            onClick={handleSearch}
          >
            Search
          </button>
        </div>
      </div>

      {/* Loading State */}
      {isLoading ? (
        <div className="flex items-center justify-center flex-grow">
          <div className="loader"></div>
          <span className="text-white ml-2">Searching {searchTerm}...</span>
        </div>
      ) : (
        <div className="flex-grow p-4">
          {/* If results exist, display them */}
          {results.length > 0 ? (
            <div className="space-y-4">
              {results.map((url, index) => (
                <p key={index} className="text-white">
                  <a
                    href={url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-400 underline"
                  >
                    {url}
                  </a>
                </p>
              ))}
            </div>
          ) : (
            <p className="text-white">No results found</p>
          )}
        </div>
      )}

      <style jsx>{`
        .loader {
          border: 4px solid rgba(255, 255, 255, 0.3);
          border-top: 4px solid #fff;
          border-radius: 50%;
          width: 20px;
          height: 20px;
          animation: spin 1s linear infinite;
        }

        @keyframes spin {
          0% {
            transform: rotate(0deg);
          }
          100% {
            transform: rotate(360deg);
          }
        }
      `}</style>
    </div>
  );
}
