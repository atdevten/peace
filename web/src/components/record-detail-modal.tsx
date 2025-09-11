"use client";

import { X, Calendar, Zap, Heart, Lock, Globe } from "lucide-react";

interface CheckIn {
  id: string;
  date: string;
  mood: 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10;
  energy: number;
  tags: string[];
  notes: string;
  habits: Record<string, boolean>;
  isPublic: boolean;
}

interface RecordDetailModalProps {
  record: CheckIn | null;
  isOpen: boolean;
  onClose: () => void;
}

function moodToEmoji(mood: CheckIn["mood"]) {
  switch (mood) {
    case 1:
      return "😣";
    case 2:
      return "😞";
    case 3:
      return "😕";
    case 4:
      return "😟";
    case 5:
      return "😐";
    case 6:
      return "🙂";
    case 7:
      return "😊";
    case 8:
      return "😄";
    case 9:
      return "😁";
    case 10:
      return "🤩";
  }
}

function getMoodLabel(mood: CheckIn["mood"]) {
  switch (mood) {
    case 1:
      return "Rough Day";
    case 2:
      return "Struggling";
    case 3:
      return "Down";
    case 4:
      return "Meh";
    case 5:
      return "Okay";
    case 6:
      return "Pretty Good";
    case 7:
      return "Good";
    case 8:
      return "Great";
    case 9:
      return "Fantastic";
    case 10:
      return "Perfect";
  }
}

function getEnergyLabel(energy: number) {
  switch (energy) {
    case 1:
      return "Need Rest";
    case 2:
      return "Low Battery";
    case 3:
      return "Getting There";
    case 4:
      return "Steady";
    case 5:
      return "Balanced";
    case 6:
      return "Energized";
    case 7:
      return "Powerful";
    case 8:
      return "Supercharged";
    case 9:
      return "Unstoppable";
    case 10:
      return "Maximum Power";
  }
}

export function RecordDetailModal({
  record,
  isOpen,
  onClose
}: RecordDetailModalProps) {
  if (!isOpen || !record) return null;

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString(undefined, {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit"
    });
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      data-oid="s7paj.j">

      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onClose}
        data-oid=".j91ro-" />


      {/* Modal */}
      <div
        className="relative w-full max-w-md bg-gradient-to-br from-gray-900 to-black border border-white/10 rounded-2xl shadow-2xl overflow-hidden"
        data-oid="jakbsfd">

        {/* Header */}
        <div
          className="flex items-center justify-between p-6 border-b border-white/10"
          data-oid="uilll:_">

          <h2
            className="text-xl font-semibold text-gray-100"
            data-oid="57.jl_s">

            Check-in Details
          </h2>
          <button
            onClick={onClose}
            className="p-2 rounded-lg hover:bg-white/10 transition-colors"
            data-oid="5e2z9.e">

            <X className="w-5 h-5 text-gray-400" data-oid="vvm:pb7" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6" data-oid="_76ueoj">
          {/* Date and Privacy */}
          <div className="flex items-center justify-between" data-oid="a29gf1m">
            <div
              className="flex items-center gap-2 text-gray-400"
              data-oid="m6w:h0t">

              <Calendar className="w-4 h-4" data-oid="iym4u_b" />
              <span className="text-sm" data-oid="iuebsxl">
                {formatDate(record.date)}
              </span>
            </div>
            <div className="flex items-center gap-2" data-oid="vymti55">
              {record.isPublic ?
              <div
                className="flex items-center gap-1 px-2 py-1 rounded-full bg-emerald-500/20 border border-emerald-500/30"
                data-oid="87m967j">

                  <Globe
                  className="w-3 h-3 text-emerald-400"
                  data-oid="d4oze-x" />


                  <span className="text-xs text-emerald-300" data-oid=":yvggig">
                    Public
                  </span>
                </div> :

              <div
                className="flex items-center gap-1 px-2 py-1 rounded-full bg-gray-500/20 border border-gray-500/30"
                data-oid="r1cl2db">

                  <Lock className="w-3 h-3 text-gray-400" data-oid="6t3n7vl" />
                  <span className="text-xs text-gray-300" data-oid="39n1si2">
                    Private
                  </span>
                </div>
              }
            </div>
          </div>

          {/* Mood Section */}
          <div className="space-y-3" data-oid="xjf5-lk">
            <h3
              className="text-sm font-medium text-gray-300 flex items-center gap-2"
              data-oid="mr4c_xz">

              <Heart className="w-4 h-4 text-red-400" data-oid="_zcjevv" />
              Mood
            </h3>
            <div
              className="flex items-center gap-3 p-4 rounded-xl bg-gradient-to-r from-red-500/10 to-emerald-500/10 border border-white/10"
              data-oid="695uf4.">

              <span className="text-3xl" data-oid="ghwg88v">
                {moodToEmoji(record.mood)}
              </span>
              <div data-oid="_3d369e">
                <div
                  className="text-lg font-semibold text-gray-100"
                  data-oid="ivs3:dr">

                  {getMoodLabel(record.mood)}
                </div>
                <div className="text-sm text-gray-400" data-oid="nzk_b8-">
                  Level {record.mood}/10
                </div>
              </div>
            </div>
          </div>

          {/* Energy Section */}
          <div className="space-y-3" data-oid=":5y3ayz">
            <h3
              className="text-sm font-medium text-gray-300 flex items-center gap-2"
              data-oid="35ldcit">

              <Zap className="w-4 h-4 text-yellow-400" data-oid="uv69mg3" />
              Energy
            </h3>
            <div
              className="flex items-center gap-3 p-4 rounded-xl bg-gradient-to-r from-red-500/10 to-emerald-500/10 border border-white/10"
              data-oid="k8ai4f.">

              <div className="text-2xl" data-oid="sox-rj:">
                {record.energy >= 8 ?
                "🚀" :
                record.energy >= 6 ?
                "⚡" :
                record.energy >= 4 ?
                "🔋" :
                "😴"}
              </div>
              <div data-oid="4d2owcp">
                <div
                  className="text-lg font-semibold text-gray-100"
                  data-oid="l3qrn6d">

                  {getEnergyLabel(record.energy)}
                </div>
                <div className="text-sm text-gray-400" data-oid="ztwqoqa">
                  Level {record.energy}/10
                </div>
              </div>
            </div>
          </div>

          {/* Tags */}
          {record.tags.length > 0 &&
          <div className="space-y-3" data-oid="mamq9v_">
              <h3
              className="text-sm font-medium text-gray-300"
              data-oid="p:sc_yt">

                Tags
              </h3>
              <div className="flex flex-wrap gap-2" data-oid="6k__m_6">
                {record.tags.map((tag) =>
              <span
                key={tag}
                className="px-3 py-1 rounded-full bg-emerald-500/20 text-emerald-300 text-sm border border-emerald-500/30"
                data-oid="m-ondit">

                    #{tag}
                  </span>
              )}
              </div>
            </div>
          }

          {/* Notes */}
          {record.notes &&
          <div className="space-y-3" data-oid="r64d-8m">
              <h3
              className="text-sm font-medium text-gray-300"
              data-oid="zk1afe_">

                Notes
              </h3>
              <div
              className="p-4 rounded-xl bg-white/[0.02] border border-white/10 max-h-48 overflow-y-auto"
              data-oid="armxkt_">

                <p
                className="text-gray-200 text-sm leading-relaxed whitespace-pre-wrap"
                data-oid="rh5tfwy">

                  {record.notes}
                </p>
              </div>
            </div>
          }

          {/* Habits */}
          <div className="space-y-3" data-oid="a7qzsl7">
            <h3
              className="text-sm font-medium text-gray-300"
              data-oid="h8p3u-c">

              Habits
            </h3>
            <div className="grid grid-cols-2 gap-2" data-oid="_ejgxn2">
              {Object.entries(record.habits).map(([habit, completed]) =>
              <div
                key={habit}
                className={`flex items-center gap-2 p-2 rounded-lg border ${
                completed ?
                "bg-emerald-500/20 border-emerald-500/30 text-emerald-300" :
                "bg-gray-500/10 border-gray-500/20 text-gray-400"}`
                }
                data-oid="b90xfyg">

                  <span className="text-sm" data-oid="yn7s6k5">
                    {completed ? "✅" : "⭕"}
                  </span>
                  <span className="text-sm capitalize" data-oid=":6dn1.l">
                    {habit}
                  </span>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>);

}