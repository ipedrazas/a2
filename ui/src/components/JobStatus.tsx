import type { JobResponse } from '../types';

interface JobStatusProps {
  job: JobResponse | null;
}

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

export function JobStatus({ job }: JobStatusProps) {
  if (!job) {
    return null;
  }

  const statusColor = statusColors[job.status];
  const statusLabel = statusLabels[job.status];

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-6 shadow-sm">
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-semibold text-gray-800">Job Status</h2>
        <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColor}`}>
          {statusLabel}
        </span>
      </div>

      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-gray-600">Job ID:</span>
          <span className="font-mono text-gray-800">{job.job_id}</span>
        </div>

        <div className="flex justify-between">
          <span className="text-gray-600">Repository:</span>
          <a
            href={job.github_url}
            target="_blank"
            rel="noopener noreferrer"
            className="text-blue-600 hover:underline"
          >
            {job.github_url}
          </a>
        </div>

        <div className="flex justify-between">
          <span className="text-gray-600">Submitted:</span>
          <span className="text-gray-800">
            {new Date(job.submitted_at).toLocaleString()}
          </span>
        </div>

        {job.started_at && (
          <div className="flex justify-between">
            <span className="text-gray-600">Started:</span>
            <span className="text-gray-800">
              {new Date(job.started_at).toLocaleString()}
            </span>
          </div>
        )}

        {job.completed_at && (
          <div className="flex justify-between">
            <span className="text-gray-600">Completed:</span>
            <span className="text-gray-800">
              {new Date(job.completed_at).toLocaleString()}
            </span>
          </div>
        )}

        {job.started_at && job.completed_at && (
          <div className="flex justify-between">
            <span className="text-gray-600">Duration:</span>
            <span className="text-gray-800">
              {Math.round(
                (new Date(job.completed_at).getTime() -
                  new Date(job.started_at).getTime()) /
                  1000
              )}{' '}
              seconds
            </span>
          </div>
        )}

        {job.error && (
          <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded">
            <p className="text-sm font-medium text-red-800">Error</p>
            <p className="text-sm text-red-700 mt-1">{job.error}</p>
          </div>
        )}
      </div>
    </div>
  );
}
