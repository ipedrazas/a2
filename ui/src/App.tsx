import { useState } from 'react';
import { BrowserRouter, Routes, Route, Link, useLocation, useNavigate } from 'react-router-dom';
import { CheckForm } from './components/CheckForm';
import { JobStatus } from './components/JobStatus';
import { Results } from './components/Results';
import { JobsList } from './components/JobsList';
import { JobDetails } from './components/JobDetails';
import { api } from './api/client';
import type { CheckRequest, JobResponse } from './types';

function Navigation() {
  const location = useLocation();
  const isActive = (path: string) => location.pathname === path;

  return (
    <nav className="bg-white border-b border-gray-200">
      <div className="max-w-5xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="flex items-center gap-2">
            <h1 className="text-xl font-bold text-gray-900">A2</h1>
          </Link>
          <div className="flex items-center gap-6">
            <Link
              to="/"
              className={`text-sm font-medium transition-colors ${
                isActive('/') ? 'text-blue-600' : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              New Check
            </Link>
            <Link
              to="/jobs"
              className={`text-sm font-medium transition-colors ${
                isActive('/jobs') || location.pathname.startsWith('/jobs/')
                  ? 'text-blue-600'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              Jobs
            </Link>
          </div>
        </div>
      </div>
    </nav>
  );
}

function HomePage() {
  const [loading, setLoading] = useState(false);
  const [currentJob, setCurrentJob] = useState<JobResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

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

      // When done, navigate to job details
      navigate(`/jobs/${jobId}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <div className="mb-8">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Run Code Quality Checks</h2>
        <p className="text-gray-600">
          Enter a GitHub repository URL to run checks and get a quality assessment.
        </p>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm mb-6">
        <CheckForm onSubmit={handleSubmit} loading={loading} />
      </div>

      {currentJob && (
        <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
          <JobStatus job={currentJob} />
          {currentJob.result && <Results result={currentJob.result} />}
        </div>
      )}
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <div className="min-h-screen bg-gray-50">
        <Navigation />

        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/jobs" element={<JobsList />} />
          <Route path="/jobs/:jobId" element={<JobDetails />} />
        </Routes>

        <footer className="border-t border-gray-200 mt-12">
          <div className="max-w-5xl mx-auto px-4 py-6 text-center text-sm text-gray-500">
            <p>A2 Code Quality Checker</p>
          </div>
        </footer>
      </div>
    </BrowserRouter>
  );
}

export default App;
