"use client";

import { useEffect, useMemo, useState } from "react";
import { ChevronDown, ChevronUp } from 'lucide-react';
import { ProtectedRoute } from '@/components/protected-route';
import { Navigation } from '@/components/navigation';
import { DailyQuote } from '@/components/daily-quote';
import { useStreak } from "@/hooks/use-streak";
import { useMentalHealthRecords, useMentalHealthHeatmap } from '@/hooks/use-mental-health-records';
import { MentalHealthRecord } from '@/lib/types';
import { BuyMeACoffee } from '@/components/buy-me-coffee';
import { RecordDetailModal } from '../components/record-detail-modal';

type HabitKey =
"meditate" |
"walk" |
"hydrate" |
"journal" |
"socialize" |
"therapy";

type CheckIn = {
  id: string;
  date: string; // ISO
  mood: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10;
  energy: number; // 1-10
  tags: string[];
  notes: string;
  habits: Record<HabitKey, boolean>;
  isPublic: boolean; // Privacy status
};

// Removed STORAGE_KEY - now using API instead of localStorage

const defaultHabits: {key: HabitKey;label: string;emoji: string;}[] = [
{ key: "meditate", label: "Meditate", emoji: "🧘" },
{ key: "walk", label: "Walk", emoji: "🚶" },
{ key: "hydrate", label: "Hydrate", emoji: "💧" },
{ key: "journal", label: "Journal", emoji: "📓" },
{ key: "socialize", label: "Socialize", emoji: "🗣️" },
{ key: "therapy", label: "Therapy", emoji: "🩺" }];


const defaultTags = [
"calm",
"grateful",
"stressed",
"overwhelmed",
"hopeful",
"anxious",
"focused",
"tired",
"content",
"lonely"];


function formatDate(d: string | Date) {
  const date = typeof d === "string" ? new Date(d) : d;
  return date.toLocaleString(undefined, {
    weekday: "short",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit"
  });
}

function MoodPicker({
  value,
  onChange
}: {value: CheckIn["mood"];onChange: (m: CheckIn["mood"]) => void;}) {
  const moods = [
    { v: 1, label: "Rough Day", emoji: "😣", color: "from-red-600 to-red-500", message: "It's okay to have tough days. Tomorrow will be better! 🌅" },
    { v: 2, label: "Struggling", emoji: "😞", color: "from-red-500 to-orange-600", message: "You're doing your best. That's what matters! 💪" },
    { v: 3, label: "Down", emoji: "😕", color: "from-orange-600 to-orange-500", message: "Small steps forward are still progress! 🚶‍♂️" },
    { v: 4, label: "Meh", emoji: "😟", color: "from-orange-500 to-yellow-500", message: "You've got this! Things will improve! ✨" },
    { v: 5, label: "Okay", emoji: "😐", color: "from-yellow-500 to-yellow-400", message: "Steady as you go! You're doing great! 🌟" },
    { v: 6, label: "Pretty Good", emoji: "🙂", color: "from-yellow-400 to-lime-500", message: "Nice! Keep that positive energy flowing! 🌊" },
    { v: 7, label: "Good", emoji: "😊", color: "from-lime-500 to-lime-400", message: "Excellent! You're shining today! ✨" },
    { v: 8, label: "Great", emoji: "😄", color: "from-lime-400 to-emerald-500", message: "Amazing! You're absolutely crushing it! 🚀" },
    { v: 9, label: "Fantastic", emoji: "😁", color: "from-emerald-500 to-emerald-400", message: "Incredible! You're unstoppable! 💫" },
    { v: 10, label: "Perfect", emoji: "🤩", color: "from-emerald-400 to-emerald-300", message: "Absolutely perfect! You're a star! ⭐" }
  ] as const;

  const selectedMood = moods.find(m => m.v === value);

  return (
    <div className="space-y-4">
      {/* Encouraging message */}
      {selectedMood && (
        <div className="text-center p-3 rounded-xl bg-gradient-to-r from-emerald-500/10 to-blue-500/10 border border-emerald-500/20">
          <p className="text-sm text-emerald-200 font-medium">{selectedMood.message}</p>
        </div>
      )}
      
      {/* Mood buttons */}
      <div className="flex items-center gap-2 flex-wrap justify-center">
        {moods.map((m) => (
          <button
            key={m.v}
            onClick={() => onChange(m.v as CheckIn["mood"])}
            className={`group relative rounded-xl p-3 md:p-4 border transition-all duration-300 transform hover:scale-110 ${
              value === m.v 
                ? "border-white/60 bg-gradient-to-br from-white/15 to-white/5 shadow-lg shadow-emerald-500/20 scale-110" 
                : "border-white/10 hover:border-white/30 hover:bg-white/5 hover:shadow-lg"
            }`}
            aria-pressed={value === m.v}
            aria-label={`${m.label} (${m.v})`}
          >
            <span
              className={`text-xl md:text-2xl transition-transform duration-300 ${
                value === m.v ? "scale-110" : "group-hover:scale-110"
              }`}
              role="img"
              aria-hidden
            >
              {m.emoji}
            </span>
            <span
              className={`absolute -bottom-2 left-1/2 -translate-x-1/2 h-2 w-8 rounded-full bg-gradient-to-r ${m.color} ${
                value === m.v ? "opacity-100" : "opacity-0 group-hover:opacity-80"
              } transition-all duration-300`}
            />
          </button>
        ))}
      </div>
    </div>
  );
}

function EnergyPicker({
  value,
  onChange
}: {value: number;onChange: (n: number) => void;}) {
  const energyLevels = [
    { v: 1, emoji: "😴", label: "Need Rest", color: "from-red-500 to-red-400", message: "Rest is productive! Take care of yourself! 🛏️" },
    { v: 2, emoji: "😪", label: "Low Battery", color: "from-red-400 to-orange-500", message: "Small energy is still energy! 🔋" },
    { v: 3, emoji: "😑", label: "Getting There", color: "from-orange-500 to-orange-400", message: "You're building momentum! 🚶‍♂️" },
    { v: 4, emoji: "😐", label: "Steady", color: "from-orange-400 to-yellow-500", message: "Consistent energy is powerful! ⚡" },
    { v: 5, emoji: "🙂", label: "Balanced", color: "from-yellow-500 to-yellow-400", message: "Perfect balance! You're in the zone! 🎯" },
    { v: 6, emoji: "😊", label: "Energized", color: "from-yellow-400 to-lime-500", message: "Great energy! Keep it flowing! 🌊" },
    { v: 7, emoji: "😄", label: "Powerful", color: "from-lime-500 to-lime-400", message: "You're unstoppable! 💪" },
    { v: 8, emoji: "😁", label: "Supercharged", color: "from-lime-400 to-emerald-500", message: "Incredible energy! You're on fire! 🔥" },
    { v: 9, emoji: "🤩", label: "Unstoppable", color: "from-emerald-500 to-emerald-400", message: "Absolutely unstoppable! 🚀" },
    { v: 10, emoji: "🚀", label: "Maximum Power", color: "from-emerald-400 to-emerald-300", message: "MAXIMUM POWER! You're incredible! ⭐" }
  ];

  const selectedEnergy = energyLevels.find(l => l.v === value);

  return (
    <div className="space-y-4">
      {/* Energy level display with message */}
      <div className="flex items-center justify-between text-sm text-gray-300">
        <span>Energy Level</span>
        <div className="flex items-center gap-2">
          <span className="text-emerald-300 font-medium">{value}/10</span>
          <span className="text-lg">{selectedEnergy?.emoji}</span>
        </div>
      </div>

      {/* Encouraging message */}
      {selectedEnergy && (
        <div className="text-center p-3 rounded-xl bg-gradient-to-r from-emerald-500/10 to-blue-500/10 border border-emerald-500/20">
          <p className="text-sm text-emerald-200 font-medium">{selectedEnergy.message}</p>
        </div>
      )}

      {/* Animated energy bars */}
      <div className="flex items-end gap-1 h-20 p-4 rounded-xl border border-white/10 bg-gradient-to-br from-white/[0.02] to-white/[0.01]">
        {energyLevels.map((level) => (
          <button
            key={level.v}
            onClick={() => onChange(level.v)}
            className={`flex-1 rounded-md transition-all duration-500 hover:scale-110 ${
              value >= level.v
                ? `bg-gradient-to-t ${level.color} opacity-100 shadow-lg`
                : "bg-white/10 opacity-40 hover:opacity-60"
            }`}
            style={{
              height: `${(level.v + 1) / 11 * 100}%`,
              minHeight: "12px"
            }}
            title={`${level.label} (${level.v})`}
          />
        ))}
      </div>

      {/* Interactive emoji selector */}
      <div className="flex items-center justify-center gap-2 flex-wrap">
        {energyLevels
          .filter((_, i) => i % 2 === 0)
          .map((level) => (
            <button
              key={level.v}
              onClick={() => onChange(level.v)}
              className={`group relative p-3 rounded-xl border transition-all duration-300 transform hover:scale-110 ${
                value === level.v
                  ? "border-emerald-400/60 bg-gradient-to-br from-emerald-500/20 to-emerald-400/10 shadow-lg shadow-emerald-500/20 scale-110"
                  : "border-white/10 hover:border-white/30 hover:bg-white/5"
              }`}
              title={level.label}
            >
              <span
                className={`text-xl transition-transform duration-300 ${
                  value === level.v ? "scale-110" : "group-hover:scale-110"
                }`}
                role="img"
                aria-hidden
              >
                {level.emoji}
              </span>
              <span
                className={`absolute -bottom-1 left-1/2 -translate-x-1/2 h-1.5 w-6 rounded-full bg-gradient-to-r ${level.color} ${
                  value === level.v ? "opacity-100" : "opacity-0 group-hover:opacity-60"
                } transition-all duration-300`}
              />
            </button>
          ))}
      </div>

      <div className="text-center text-xs text-gray-400">
        {selectedEnergy?.label || "Select your energy level"}
      </div>
    </div>
  );
}

// Removed unused Slider component

function TagPicker({
  selected,
  onToggle



}: {selected: string[];onToggle: (t: string) => void;}) {
  return (
    <div className="flex flex-wrap gap-2" data-oid="h4ywees">
      {defaultTags.map((t) => {
        const active = selected.includes(t);
        return (
          <button
            key={t}
            onClick={() => onToggle(t)}
            className={`px-3 py-1.5 rounded-full text-sm transition border ${
            active ?
            "bg-emerald-500/20 text-emerald-200 border-emerald-400/30" :
            "bg-white/5 text-gray-300 border-white/10 hover:border-white/25 hover:bg-white/10"}`
            }
            data-oid=":-x9e83">

            #{t}
          </button>);

      })}
    </div>);

}

function HabitList({
  habits,
  onToggle



}: {habits: Record<HabitKey, boolean>;onToggle: (k: HabitKey) => void;}) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 gap-2" data-oid="m8o8.k9">
      {defaultHabits.map((h) => {
        const active = habits[h.key];
        return (
          <button
            key={h.key}
            onClick={() => onToggle(h.key)}
            className={`flex items-center gap-2 rounded-xl px-3 py-2 border text-sm transition ${
            active ?
            "bg-emerald-500/15 text-emerald-100 border-emerald-400/30" :
            "bg-white/5 text-gray-300 border-white/10 hover:border-white/25 hover:bg-white/10"}`
            }
            data-oid="h.af551">

            <span className="text-base" aria-hidden data-oid="n5rwdk5">
              {h.emoji}
            </span>
            <span data-oid="b_39n:h">{h.label}</span>
          </button>);

      })}
    </div>);

}

function PrivacyToggle({
  isPublic,
  onChange



}: {isPublic: boolean;onChange: (value: boolean) => void;}) {
  return (
    <div
      className="flex items-center justify-between p-3 rounded-xl border border-white/10 bg-white/[0.02]"
      data-oid="yssljt6">

      <div className="flex items-center gap-3" data-oid="m8k550_">
        <span className="text-lg" role="img" aria-hidden data-oid="upd386v">
          {isPublic ? "🌍" : "🔒"}
        </span>
        <div data-oid="4k92_kq">
          <div className="text-sm font-medium text-gray-200" data-oid="lesiznj">
            {isPublic ? "Public Entry" : "Private Entry"}
          </div>
          <div className="text-xs text-gray-400" data-oid="w48lwxv">
            {isPublic ?
            "This entry can be shared with others" :
            "This entry is only visible to you"}
          </div>
        </div>
      </div>
      <button
        onClick={() => onChange(!isPublic)}
        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
        isPublic ? "bg-emerald-500" : "bg-gray-600"}`
        }
        role="switch"
        aria-checked={isPublic}
        data-oid="4bimttt">

        <span
          className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
          isPublic ? "translate-x-6" : "translate-x-1"}`
          }
          data-oid="z:v8zv7" />

      </button>
    </div>);

}

function Card({
  title,
  subtitle,
  children,
  right





}: {title: string;subtitle?: string;children: React.ReactNode;right?: React.ReactNode;}) {
  return (
    <section
      className="rounded-2xl border border-white/10 bg-white/[0.04] backdrop-blur p-4 md:p-5 shadow-lg shadow-black/20"
      data-oid="8rn5akj">

      <div className="mb-3 flex items-start justify-between" data-oid="h.xmnrz">
        <div data-oid="wp4dwi3">
          <h3
            className="text-gray-100 font-medium tracking-tight"
            data-oid="ecmt9hc">

            {title}
          </h3>
          {subtitle ?
          <p className="text-sm text-gray-400 mt-0.5" data-oid="vuon-a0">
              {subtitle}
            </p> :
          null}
        </div>
        {right}
      </div>
      {children}
    </section>);

}

// Removed useLocalCheckins - now using API via useMentalHealthRecords hook

function useStats(history: CheckIn[]) {
  return useMemo(() => {
    if (!history.length) return { avg7: 0, entries7: 0 };

    const byDay = new Map<string, CheckIn[]>();
    for (const h of history) {
      const d = new Date(h.date);
      const key = d.toISOString().slice(0, 10);
      const arr = byDay.get(key) || [];
      arr.push(h);
      byDay.set(key, arr);
    }

    // last 7 days
    const today = new Date();
    const days: string[] = [];
    for (let i = 0; i < 7; i++) {
      const d = new Date(today);
      d.setDate(today.getDate() - i);
      days.push(d.toISOString().slice(0, 10));
    }

    let moodSum = 0;
    let count = 0;
    for (const day of days) {
      const entries = byDay.get(day) || [];
      if (entries.length) {
        count++;
        moodSum += Math.round(
          entries.reduce((s, e) => s + e.mood, 0) / entries.length
        );
      }
    }

    return {
      avg7: count ? Math.round(moodSum / count * 10) / 10 : 0,
      entries7: count
    };  }, [history]);
}
export default function Page() {
  const { records, loading: recordsLoading, createRecord } = useMentalHealthRecords();
  
  // Convert API records to local CheckIn format for compatibility with existing UI
  const history = useMemo(() => {
    if (!records || !Array.isArray(records)) {
      return [];
    }
    
    return records.map((record: MentalHealthRecord): CheckIn => ({
      id: record.id,
      date: record.created_at,
      mood: record.happy_level as CheckIn["mood"],
      energy: record.energy_level,
      tags: [], // Ignore tags for now as requested
      notes: record.notes || "",
      habits: {
        meditate: false,
        walk: false,
        hydrate: false,
        journal: false,
        socialize: false,
        therapy: false
      }, // Default habits since backend doesn't have this yet
      isPublic: record.status === "public"
    }));
  }, [records]);
  
  // Get streak from server API
  const { streak, loading: streakLoading, refetch: refetchStreak } = useStreak();
  const stats = useStats(history);

  // Check if user has already checked in today (Option 3: Smart state)
  const todayEntry = useMemo(() => {
    const today = new Date().toISOString().slice(0, 10);
    return history.find(entry => 
      new Date(entry.date).toISOString().slice(0, 10) === today
    );
  }, [history]);

  // Option 3: Default to expanded if no today entry, collapsed if already checked in
  // Option 1: Allow manual toggle
  const [isCheckInExpanded, setIsCheckInExpanded] = useState<boolean>(!todayEntry);

  // Update expansion state when todayEntry changes (when records load)
  useEffect(() => {
    setIsCheckInExpanded(!todayEntry);
  }, [todayEntry]);

  const [mood, setMood] = useState<CheckIn["mood"]>(5);
  const [energy, setEnergy] = useState(5);
  const [tags, setTags] = useState<string[]>([]);
  const [notes, setNotes] = useState("");
  const [isPublic, setIsPublic] = useState(false);
  const [habits, setHabits] = useState<Record<HabitKey, boolean>>({
    meditate: false,
    walk: false,
    hydrate: false,
    journal: false,
    socialize: false,
    therapy: false
  });
  
  // Loading and error states for save operation
  const [isSaving, setIsSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);
  
  // Modal state for record details
  const [selectedRecord, setSelectedRecord] = useState<CheckIn | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const toggleTag = (t: string) =>
  setTags((prev) =>
  prev.includes(t) ? prev.filter((x) => x !== t) : [...prev, t]
  );
  const toggleHabit = (k: HabitKey) =>
  setHabits((prev) => ({ ...prev, [k]: !prev[k] }));

  const handleRecordClick = (record: CheckIn) => {
    setSelectedRecord(record);
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setSelectedRecord(null);
  };

  function resetForm() {
    setMood(5);
    setEnergy(5);
    setTags([]);
    setNotes("");
    setIsPublic(false);
    setHabits({
      meditate: false,
      walk: false,
      hydrate: false,
      journal: false,
      socialize: false,
      therapy: false
    });
  }

  async function saveCheckIn() {
    setIsSaving(true);
    setSaveError(null);
    
    try {
      // Map frontend data to backend API format
      const recordData = {
        happy_level: mood, // mood (1-10) maps directly to happy_level
        energy_level: energy, // energy (1-10) maps directly to energy_level
        notes: notes.trim() || undefined, // Send undefined if empty to make it null on backend
        status: isPublic ? "public" : "private"
      };
      
      await createRecord(recordData);
      // Refresh streak after creating new record
      refetchStreak();
      resetForm();
      // Auto-collapse after successful save (Option 1 + smart behavior)
      setIsCheckInExpanded(false);
      // Refetch to ensure UI is in sync (createRecord should update the list automatically)
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to save check-in';
      setSaveError(errorMessage);
    } finally {
      setIsSaving(false);
    }
  }

  const last5 = history.slice(0, 5);

  return (
    <ProtectedRoute>
      <div className="min-h-screen w-full bg-[radial-gradient(1200px_600px_at_100%_-10%,rgba(16,185,129,0.08),transparent_60%),radial-gradient(800px_400px_at_0%_0%,rgba(59,130,246,0.08),transparent_50%)] bg-black text-gray-100 relative">
        <Navigation />
        <main className="mx-auto max-w-6xl px-4 py-10 md:py-14">
          {/* Header */}
          <header className="mb-8 md:mb-10">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h1 className="text-3xl md:text-4xl font-semibold tracking-tight">
                  Mindful — Track, Care, Grow
                </h1>
              </div>
              <div className="hidden md:flex items-center gap-2 text-sm text-gray-400">
                <div className="rounded-xl border border-white/10 px-3 py-2 bg-white/[0.03]">
                  <span className="text-gray-300">7‑day avg</span>
                  <span className="ml-2 text-emerald-300 font-medium">
                    {stats.avg7 || "–"}
                  </span>
                </div>
                <div className="rounded-xl border border-white/10 px-3 py-2 bg-white/[0.03]">
                  <span className="text-gray-300">Streak</span>
                  <span className="ml-2 text-emerald-300 font-medium">
                    {streakLoading ? "..." : streak}d
                  </span>
                </div>
                <BuyMeACoffee username="svie4mv" compact={true} />
              </div>
            </div>
          </header>

          {/* Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-5">
            {/* Left: Check-in */}
            <div className="lg:col-span-2 space-y-4 md:space-y-5">
              <Card
                title="Today's check‑in"
                subtitle={
                  todayEntry 
                    ? `Already checked in today (Mood: ${todayEntry.mood}/10, Energy: ${todayEntry.energy}/10)` 
                    : "How are you feeling right now?"
                }
                right={
                  <button
                    onClick={() => setIsCheckInExpanded(!isCheckInExpanded)}
                    className="p-2 rounded-lg border border-white/10 hover:bg-white/10 hover:border-white/30 transition-all"
                    aria-label={isCheckInExpanded ? "Collapse check-in form" : "Expand check-in form"}
                  >
                    {isCheckInExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                  </button>
                }
              >
                {isCheckInExpanded && (
                  <div className="flex flex-col gap-4">
                    <MoodPicker value={mood} onChange={setMood} />
                    <EnergyPicker value={energy} onChange={setEnergy} />
                    <div className="flex flex-col gap-2">
                      <span className="text-sm text-gray-300">Tags</span>
                      <TagPicker selected={tags} onToggle={toggleTag} />
                    </div>
                    <div className="flex flex-col gap-2">
                      <label htmlFor="notes" className="text-sm text-gray-300">Notes</label>
                      <textarea
                        id="notes"
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        placeholder="Optional journal…"
                        rows={3}
                        className="w-full resize-none rounded-xl border border-white/10 bg-black/30 px-3 py-2 text-gray-100 placeholder:text-gray-500 focus:border-emerald-400/40 focus:outline-none"
                      />
                    </div>
                    <div className="flex flex-col gap-2">
                      <span className="text-sm text-gray-300">Privacy</span>
                      <PrivacyToggle isPublic={isPublic} onChange={setIsPublic} />
                    </div>
                    {saveError && (
                      <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/20">
                        <p className="text-red-400 text-sm">{saveError}</p>
                      </div>
                    )}
                    <div className="flex items-center justify-between gap-3">
                      <div className="text-sm text-gray-400">
                        Synced with server • {new Date().toLocaleDateString()}
                        {recordsLoading && (
                          <span className="ml-2 inline-flex items-center gap-1">
                            <div className="animate-spin rounded-full h-3 w-3 border-b border-emerald-400"></div>
                            Loading...
                          </span>
                        )}
                      </div>
                      <div className="flex items-center gap-2">
                        <button
                          onClick={resetForm}
                          disabled={isSaving}
                          className="rounded-xl border border-white/10 bg-white/[0.02] px-4 py-2 text-gray-300 hover:bg-white/10 hover:border-white/20 transition disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          Reset
                        </button>
                        <button
                          onClick={saveCheckIn}
                          disabled={isSaving}
                          className="rounded-xl border border-emerald-500/30 bg-emerald-500/15 px-4 py-2 text-emerald-100 hover:bg-emerald-500/25 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                        >
                          {isSaving && (
                            <div className="animate-spin rounded-full h-4 w-4 border-b border-emerald-400"></div>
                          )}
                          {isSaving ? 'Saving...' : 'Save check‑in'}
                        </button>
                      </div>
                    </div>
                  </div>
                )}
              </Card>

              {/* OPTION 1: Daily Quote Full Width Below Check-in */}
              <DailyQuote />

              <Card title="Habits" subtitle="Small actions that support your mood">
                <HabitList habits={habits} onToggle={toggleHabit} />
              </Card>
            </div>

            {/* Right: Insights + History */}
            <div className="space-y-4 md:space-y-5">
              <Card
                title="Insights"
                subtitle="Last 7 days overview"
                right={
                  <span className="text-xs text-gray-500">
                    {stats.entries7}d logged
                  </span>
                }
              >
                <div className="grid grid-cols-3 gap-3">
                  <div className="rounded-xl border border-white/10 bg-white/[0.02] p-3">
                    <div className="text-xs text-gray-400">Avg mood</div>
                    <div className="mt-1 text-xl font-medium text-emerald-300">
                      {stats.avg7 || "–"}
                    </div>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-white/[0.02] p-3">
                    <div className="text-xs text-gray-400">Streak</div>
                    <div className="mt-1 text-xl font-medium text-emerald-300">
                      {streakLoading ? "..." : streak}d
                    </div>
                  </div>
                  <div className="rounded-xl border border-white/10 bg-white/[0.02] p-3">
                    <div className="text-xs text-gray-400">Entries</div>
                    <div className="mt-1 text-xl font-medium text-emerald-300">
                      {history.length}
                    </div>
                  </div>
                </div>
                <div className="mt-4">
                  <MiniBars history={history} />
                </div>
                <div className="mt-6">
                  <MoodHeatmap />
                </div>
              </Card>

              <Card title="Recent" subtitle="Your last entries">
                {recordsLoading ? (
                  <div className="flex items-center justify-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-emerald-500"></div>
                  </div>
                ) : last5.length === 0 ? (
                  <p className="text-sm text-gray-400">
                    No entries yet. Log your first check‑in above.
                  </p>
                ) : (
                  <ul className="flex flex-col divide-y divide-white/5">
                    {last5.map((e) => (
                      <li 
                        key={e.id} 
                        className="py-3 cursor-pointer hover:bg-white/[0.02] rounded-lg transition-colors"
                        onClick={() => handleRecordClick(e)}
                      >
                        <div className="flex items-start justify-between gap-3 px-2">
                          <div className="flex items-center gap-3">
                            <span className="text-xl" aria-hidden>
                              {moodToEmoji(e.mood)}
                            </span>
                            <div>
                              <div className="flex items-center gap-2">
                                <div className="text-sm text-gray-200">
                                  {e.tags.length ? e.tags.map((t) => `#${t}`).join("  ") : "No tags"}
                                </div>
                                <span
                                  className="text-xs"
                                  role="img"
                                  aria-label={e.isPublic ? "Public" : "Private"}
                                >
                                  {e.isPublic ? "🌍" : "🔒"}
                                </span>
                              </div>
                              {e.notes && (
                                <div className="mt-1 text-sm text-gray-400 line-clamp-2 max-w-[36ch]">
                                  {e.notes}
                                </div>
                              )}
                              <div className="mt-1 flex gap-2 text-[10px] text-gray-500">
                                <span>Energy {e.energy}</span>
                              </div>
                            </div>
                          </div>
                          <div className="text-xs text-gray-500 whitespace-nowrap">
                            {formatDate(e.date)}
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>
                )}
              </Card>
            </div>
          </div>

          <footer className="mt-10 text-center text-xs text-gray-500">
            Built with love. If you need help, talk to someone you trust. This is not medical advice.
          </footer>
        </main>
      </div>
      
      {/* Record Detail Modal */}
      <RecordDetailModal
        record={selectedRecord}
        isOpen={isModalOpen}
        onClose={handleCloseModal}
      />
    </ProtectedRoute>
  );
}

function moodToEmoji(mood: CheckIn["mood"]) {
  switch (mood) {
    case 1: return "😣";
    case 2: return "😞";
    case 3: return "😕";
    case 4: return "😟";
    case 5: return "😐";
    case 6: return "🙂";
    case 7: return "😊";
    case 8: return "😄";
    case 9: return "😁";
    case 10: return "🤩";
  }
}

function MiniBars({ history }: {history: CheckIn[];}) {
  const today = new Date();
  const days = Array.from({ length: 7 }, (_, i) => {
    const d = new Date(today);
    d.setDate(today.getDate() - (6 - i));
    return d.toISOString().slice(0, 10);
  });

  const byDay: Record<string, number | null> = {};
  for (const day of days) byDay[day] = null;

  for (const e of history) {
    const key = new Date(e.date).toISOString().slice(0, 10);
    if (key in byDay) {
      byDay[key] = byDay[key] == null ? e.mood : Math.round(((byDay[key] as number) + e.mood) / 2);
    }
  }

  return (
    <div className="grid grid-cols-7 gap-2 items-end">
      {days.map((d) => {
        const v = byDay[d];
        const h = v ? v / 10 * 64 + 8 : 8; // px height
        const color = v ? v >= 8 ? "bg-emerald-400" : v >= 6 ? "bg-lime-400" : v >= 4 ? "bg-yellow-400" : "bg-orange-400" : "bg-white/10";
        return (
          <div key={d} className="flex flex-col items-center gap-1">
            <div className={`w-6 rounded-md ${color}`} style={{ height: h }} />
            <div className="text-[10px] text-gray-500">
              {new Date(d).toLocaleDateString(undefined, { weekday: "short" }).slice(0, 2)}
            </div>
          </div>
        );
      })}
    </div>
  );
}

function MoodHeatmap() {
  // Use the heatmap API hook instead of local history
  const today = new Date();
  const startDate = new Date(today);
  startDate.setDate(today.getDate() - 90); // Show last 90 days
  
  // Set to start of start day (00:00:00) and start of next day after today in RFC3339 format
  const startedAt = new Date(startDate.getFullYear(), startDate.getMonth(), startDate.getDate()).toISOString();
  const nextDay = new Date(today.getFullYear(), today.getMonth(), today.getDate() + 1);
  const endedAt = nextDay.toISOString();
  
  const { heatmapData, loading, error } = useMentalHealthHeatmap(startedAt, endedAt);

  // Generate all days in the range
  const days: string[] = [];
  const current = new Date(startDate);
  
  // Generate days for the full range
  while (current <= today) {
    days.push(current.toISOString().slice(0, 10));
    current.setDate(current.getDate() + 1);
  }
  
  // Also add any dates from API data that might be outside our range
  if (heatmapData?.data) {
    Object.keys(heatmapData.data).forEach(apiDate => {
      if (!days.includes(apiDate)) {
        days.push(apiDate);
      }
    });
    // Sort the days to maintain chronological order
    days.sort();
  }

  // Convert API data to display format
  const byDay: Record<string, number | null> = {};
  for (const day of days) byDay[day] = null;

  if (heatmapData?.data) {
    Object.entries(heatmapData.data).forEach(([date, dataPoint]) => {
      if (date in byDay) {
        // Use happy_level from the API response
        byDay[date] = dataPoint.happy_level;
      }
    });
  }

  // Group days by weeks
  const weeks: string[][] = [];
  let currentWeek: string[] = [];

  for (let i = 0; i < days.length; i++) {
    const day = days[i];
    const date = new Date(day);
    const dayOfWeek = date.getDay(); // 0 = Sunday

    if (i === 0) {
      // Fill empty days at the start of first week
      for (let j = 0; j < dayOfWeek; j++) {
        currentWeek.push("");
      }
    }

    currentWeek.push(day);

    if (dayOfWeek === 6 || i === days.length - 1) {
      // Saturday or last day
      // Fill empty days at the end of last week
      while (currentWeek.length < 7) {
        currentWeek.push("");
      }
      weeks.push(currentWeek);
      currentWeek = [];
    }
  }

  const getMoodColor = (mood: number | null) => {
    if (mood === null) return "bg-white/5";
    if (mood >= 9) return "bg-emerald-500";
    if (mood >= 7) return "bg-emerald-400";
    if (mood >= 5) return "bg-yellow-400";
    if (mood >= 3) return "bg-orange-400";
    return "bg-red-400";
  };

  const getMoodIntensity = (mood: number | null) => {
    if (mood === null) return "opacity-20";
    if (mood >= 8) return "opacity-100";
    if (mood >= 6) return "opacity-80";
    if (mood >= 4) return "opacity-60";
    return "opacity-40";
  };

  const weekdays = ["S", "M", "T", "W", "T", "F", "S"];
  const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

  if (loading) {
    return (
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-sm font-medium text-gray-200">Mood Heatmap</h4>
        </div>
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-emerald-500"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <h4 className="text-sm font-medium text-gray-200">Mood Heatmap</h4>
        </div>
        <div className="p-3 rounded-xl bg-red-500/10 border border-red-500/20">
          <p className="text-red-400 text-xs">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-medium text-gray-200">Mood Heatmap</h4>
        <div className="flex items-center gap-1 text-xs text-gray-400">
          <span>Less</span>
          <div className="flex gap-1">
            <div className="w-2 h-2 rounded-sm bg-white/5"></div>
            <div className="w-2 h-2 rounded-sm bg-red-400 opacity-40"></div>
            <div className="w-2 h-2 rounded-sm bg-orange-400 opacity-60"></div>
            <div className="w-2 h-2 rounded-sm bg-yellow-400 opacity-80"></div>
            <div className="w-2 h-2 rounded-sm bg-emerald-400 opacity-100"></div>
          </div>
          <span>More</span>
        </div>
      </div>

      <div className="overflow-x-auto">
        <div className="flex gap-1 min-w-fit">
          {/* Weekday labels */}
          <div className="flex flex-col gap-1 mr-2">
            <div className="h-3"></div> {/* Space for month labels */}
            {weekdays.map((day, i) => (
              <div key={i} className="w-3 h-3 flex items-center justify-center text-[10px] text-gray-500">
                {i % 2 === 1 ? day : ""}
              </div>
            ))}
          </div>

          {/* Heatmap grid */}
          <div className="flex flex-col gap-1">
            {/* Month labels */}
            <div className="flex gap-1 h-3">
              {weeks.map((week, weekIndex) => {
                const firstDay = week.find((day) => day !== "");
                if (!firstDay) return <div key={weekIndex} className="w-3"></div>;

                const date = new Date(firstDay);
                const isFirstWeekOfMonth = date.getDate() <= 7;

                return (
                  <div key={weekIndex} className="w-3 text-[10px] text-gray-500">
                    {isFirstWeekOfMonth ? months[date.getMonth()].slice(0, 3) : ""}
                  </div>
                );
              })}
            </div>

            {/* Days grid */}
            <div className="flex gap-1">
              {weeks.map((week, weekIndex) => (
                <div key={weekIndex} className="flex flex-col gap-1">
                  {week.map((day, dayIndex) => {
                    if (!day) return <div key={dayIndex} className="w-3 h-3"></div>;

                    const mood = byDay[day];
                    const date = new Date(day);
                    const isToday = day === today.toISOString().slice(0, 10);

                    return (
                      <div
                        key={day}
                        className={`w-3 h-3 rounded-sm ${getMoodColor(mood)} ${getMoodIntensity(mood)} ${
                          isToday ? "ring-1 ring-white/60" : ""
                        } transition-all hover:scale-110 cursor-pointer`}
                        title={`${date.toLocaleDateString()}: ${mood ? `Mood ${mood}` : "No entry"}`}
                      />
                    );
                  })}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="text-xs text-gray-500">
        {heatmapData?.total_records || 0} entries in the last 90 days
      </div>
    </div>
  );
}
