import { useQuery } from "@tanstack/react-query";
import { Patient } from "../types/Patient";
import axios from "axios";
import { Entity } from "../types/Entity";
import { AreaChart, ResponsiveContainer } from "recharts";
import { Area } from "recharts";
import { AiTwotoneHeart } from "react-icons/ai";
import { Scheme } from "../scheme";

interface Props {
  patient: Patient;
  height: number;
  width: number;
}

export const HeartbeatChart: React.FC<Props> = ({ patient, width, height }) => {
  const { data, error, isFetching } = useQuery<Entity>({
    queryKey: [`heartbeat.${patient.id}`],
    queryFn: async () =>
      axios
        .get(
          `http://localhost:8080/user/${patient.id}/entities/sensor.heartbeat/`
        )
        .then((res) => res.data),
  });

  if (isFetching) {
    return <h1>loading...</h1>;
  }

  return (
    <div className="p-4 grid grid-rows-2">
      <div className="row-start-1 row-span-1 flex flex-row justify-between">
        <div>
          <h1 className="text-2xl text-white">
            Pulsen ser{" "}
            <span className="text-secondary-text font-bold">fin</span> ut
          </h1>
        </div>
        <div className="animate-pulse">
          <AiTwotoneHeart size={48} color={Scheme.Red} />
        </div>
      </div>
      <div>
        <ResponsiveContainer height={160}>
          <AreaChart
            data={data.history}
            margin={{
              top: 10,
              right: 30,
              left: 0,
              bottom: 0,
            }}
          >
            <Area
              type="monotone"
              dataKey="state"
              stroke="#8884d8"
              fill="#8884d8"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
