"use client";

import { useState, useEffect } from 'react';
import { apiService } from '@/lib/api';
import { MentalHealthRecord, CreateMentalHealthRecordRequest, MentalHealthHeatmapResponse } from '@/lib/types';

export function useMentalHealthRecords(startedAt?: string, endedAt?: string) {
  const [records, setRecords] = useState<MentalHealthRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRecords = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await apiService.getMentalHealthRecords(startedAt, endedAt);
      setRecords(Array.isArray(data) ? data : []);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch records');
      setRecords([]); // Ensure records is always an array, even on error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecords();
  }, [startedAt, endedAt]);

  const createRecord = async (data: CreateMentalHealthRecordRequest) => {
    try {
      const newRecord = await apiService.createMentalHealthRecord(data);
      if (newRecord) {
        setRecords(prev => Array.isArray(prev) ? [newRecord, ...prev] : [newRecord]);
      }
      return newRecord;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to create record');
    }
  };

  const deleteRecord = async (id: string) => {
    try {
      await apiService.deleteMentalHealthRecord(id);
      setRecords(prev => Array.isArray(prev) ? prev.filter(record => record.id !== id) : []);
    } catch (err: any) {
      throw new Error(err.message || 'Failed to delete record');
    }
  };

  return {
    records,
    loading,
    error,
    createRecord,
    deleteRecord,
    refetch: fetchRecords,
  };
}

export function useMentalHealthHeatmap(startedAt?: string, endedAt?: string) {
  const [heatmapData, setHeatmapData] = useState<MentalHealthHeatmapResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchHeatmap = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await apiService.getMentalHealthHeatmap(startedAt, endedAt);
        setHeatmapData(data);
      } catch (err: any) {
        setError(err.message || 'Failed to fetch heatmap data');
      } finally {
        setLoading(false);
      }
    };

    fetchHeatmap();
  }, [startedAt, endedAt]);

  return {
    heatmapData,
    loading,
    error,
  };
}
