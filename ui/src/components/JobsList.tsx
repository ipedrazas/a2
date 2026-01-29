import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import type { JobResponse } from '../types';
import { api } from '../api/client';

const statusColors = {
  pending: 'bg-yellow-100 text-yellow-800',
  running: 'bg-blue-100 text-blue-800',
  completed: 'bg-green-100 text-green-800',
  failed: 'bg-red-100 text-red-800',
};

const statusLabels = {
  pending: 'Pending',
  running: 'Running',
  completed: 'Completed',
  failed: 'Failed',
};

export function JobsList() {
  const [jobs, setJobs] = useState<JobResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    loadJobs();
    // Refresh every 5 seconds
    const interval = setInterval(loadJobs, 5000);
    return () => clearInterval(interval);
  }, []);

  async function loadJobs() {
    try {
      setLoading(true);
      const data = await api.listJobs();
      setJobs(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load jobs');
    } finally {
      setLoading(false);
    }
  }

  function handleJobClick(jobId: string) {
    navigate(`/jobs/${jobId}`);
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Jobs</h1>
        <button
          onClick={loadJobs}
          disabled={loading}
          className="px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-gray-400 transition-colors"
        >
          {loading ? 'Loading...' : 'Refresh'}
        </button>
      </div>

      {error && (
        <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
          <p className="text-red-800">{error}</p>
        </div>
      )}

      {jobs.length === 0 && !loading && (
        <div className="text-center py-12 bg-white border border-gray-200 rounded-lg">
          <p className="text-gray-500">No jobs yet. Submit a check to get started!</p>
        </div>
      )}

      <div className="space-y-4">
        {jobs.map((job) => (
          <div
            key={job.job_id}
            onClick={() => handleJobClick(job.job_id)}
            className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm hover:shadow-md hover:border-blue-300 cursor-pointer transition-all"
          >
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2">
                  <h3 className="text-lg font-semibold text-gray-800">{job.github_url}</h3>
                  <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColors[job.status]}`}>
                    {statusLabels[job.status]}
                  </span>
                </div>
                <div className="text-sm text-gray-600 space-y-1">
                  <p>
                    <span className="font-medium">Job ID:</span> {job.job_id}
                  </p>
                  <p>
                    <span className="font-medium">Submitted:</span>{' '}
                    {new Date(job.submitted_at).toLocaleString()}
                  </p>
                  {job.started_at && (
                    <p>
                      <span className="font-medium">Started:</span>{' '}
                      {new Date(job.started_at).toLocaleString()}
                    </p>
                  )}
                  {job.completed_at && (
                    <p>
                      <span className="font-medium">Completed:</span>{' '}
                      {new Date(job.completed_at).toLocaleString()}
                    </p>
                  )}
                  {job.result && (
                    <p>
                      <span className="font-medium">Score:</span>{' '}
                      <span className={`font-bold ${
                        job.result.summary.score >= 80 ? 'text-green-600' :
                        job.result.summary.score >= 60 ? 'text-yellow-600' : 'text-red-600'
                      }`}>
                        {job.result.summary.score.toFixed(1)}%
                      </span>
                    </p>
                  )}
                </div>
              </div>
              <div className="text-blue-600">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
