"use client";

import { useState } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import axios from "axios"
import https from "https";

export default function Home() {
  const [isLoading, setIsLoading] = useState(false); //State for loading status
  const [website, setWebsite] = useState(""); //To store the website being crawled
  const [results, setResults] = useState([]); //State to hold search results

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const agent = new https.Agent({
    rejectUnauthorized: false,
  });

  const handleSearch = async () => {
    setIsLoading(true);
    console.log("handleSearch called");
    try {
      const response = await axios.post(
        apiUrl,
        { website: website },
        { httpsAgent: agent }
      );
      setResults(response.data);
    } catch (error) {
      console.error("Error calling API:", error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-black">
      <h1 className="text-4xl font-bold text-white mb-4">SiteSearch</h1>
      {!isLoading ? (
        <div className="flex gap-2">
          <input
            type="text"
            className="border border-gray-300 rounded-md p-2 w-64 text-black"
            placeholder="Search..."
            value={website}
            onChange={(e) => setWebsite(e.target.value)}
          />
          <button
            className="bg-blue-500 text-white rounded px-4 h-full transition duration-300 hover:bg-blue-600"
            onClick={handleSearch}
          >
            Search
          </button>
        </div>
      ) : (
        <div className="flex items-center">
          <div className="loader"></div>
          <span className="text-white ml-2">Crawling {website}...</span>
        </div>
      )}
      <ScrollArea className="h-[200px] w-[350px] rounded-md border p-4">
        {results.length > 0 ? (
          results.map((url, index) => (
            <p key={index}>
              <a
                href={url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-400 underline"
              >
                {url}
              </a>
            </p>
          ))
        ) : (
          <p className="text-white">No results found</p>
        )}
      </ScrollArea>
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
