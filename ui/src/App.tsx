import { useState } from 'react';
import { CheckForm } from './components/CheckForm';
import { JobStatus } from './components/JobStatus';
import { Results } from './components/Results';
import { api } from './api/client';
import type { CheckRequest, JobResponse } from './types';

function App() {
  const [loading, setLoading] = useState(false);
  const [currentJob, setCurrentJob] = useState<JobResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (request: CheckRequest) => {
    setLoading(true);
    setError(null);
    setCurrentJob(null);

    try {
      // Submit the job
      const response = await api.submitCheck(request);
      const jobId = response.job_id;

      // Poll for results
      await api.pollJob(
        jobId,
        (job) => {
          setCurrentJob(job);
        },
        2000
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setCurrentJob(null);
    setError(null);
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-5xl mx-auto px-4 py-6">
          <h1 className="text-3xl font-bold text-gray-900">A2</h1>
          <p className="text-gray-600 mt-1">
            Code Quality Checker - Run checks on any GitHub repository
          </p>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-5xl mx-auto px-4 py-8">
        {error && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-red-800">{error}</p>
          </div>
        )}

        {!currentJob ? (
          <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
            <h2 className="text-xl font-semibold text-gray-800 mb-4">
              Run Code Quality Checks
            </h2>
            <CheckForm onSubmit={handleSubmit} loading={loading} />
          </div>
        ) : (
          <div className="space-y-6">
            <JobStatus job={currentJob} />

            {(currentJob.status === 'completed' || currentJob.status === 'failed') && (
              <div className="flex justify-center">
                <button
                  onClick={handleReset}
                  className="px-6 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Run Another Check
                </button>
              </div>
            )}

            {currentJob.result && <Results result={currentJob.result} />}
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-200 mt-12">
        <div className="max-w-5xl mx-auto px-4 py-6 text-center text-sm text-gray-500">
          <p>A2 Code Quality Checker</p>
        </div>
      </footer>
    </div>
  );
}

export default App;
