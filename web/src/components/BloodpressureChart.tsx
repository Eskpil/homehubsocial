import { Patient } from "../types/Patient";
import { Scheme } from "../scheme";
import { AiTwotoneHeart } from "react-icons/ai";
import { ResponsiveContainer, AreaChart, Area } from "recharts";
import { useQueries, useQuery } from "@tanstack/react-query";
import { Entity } from "../types/Entity";
import axios from "axios";
import { useState } from "react";

interface Props {
  patient: Patient;
}

interface Group {
  systolic: string;
  diastolic: string;
}

export const BloodpressureGraph: React.FC<Props> = ({ patient }) => {
  //const {systolicData, systolicError, systolicIsFetching, systolicIsSuccess} = useQuery<Entity>({
  //    queryKey: [`systolic.${patient.id}`],
  //    queryFn: async () => axios.get(`http://localhost:8080/user/${patient.id}/entities/sensor.bloodpressure_systolic/`).then((res) => res.data)
  //})

  //const {diastolicData, diastolicError, diastolicIsFetching, diastolicIsSuccess} = useQuery<Entity>({
  //    queryKey: [`diastolic.${patient.id}`],
  //    queryFn: async () => axios.get(`http://localhost:8080/user/${patient.id}/entities/sensor.bloodpressure_diastolic/`).then((res) => res.data)
  //})

  const [systolicQuery, diastolicQuery] = useQueries<Entity[]>({
    queries: [
      {
        queryKey: [`systolic`, patient.id],
        queryFn: async () =>
          axios
            .get(
              `http://localhost:8080/user/${patient.id}/entities/sensor.bloodpressure_systolic/`
            )
            .then((res) => res.data),
      },
      {
        queryKey: [`diastolic`, patient.id],
        queryFn: async () =>
          axios
            .get(
              `http://localhost:8080/user/${patient.id}/entities/sensor.bloodpressure_diastolic/`
            )
            .then((res) => res.data),
      },
    ],
  });

  if (systolicQuery.isLoading) return <div>loading...</div>;
  if (diastolicQuery.isLoading) return <div>loading...</div>;

  if (systolicQuery.isError) return <div>error...</div>;
  if (diastolicQuery.isError) return <div>error...</div>;

  const data: Group[] = [];

  (systolicQuery.data as Entity).history.forEach((systolic, idx) => {
    const diastolic = (diastolicQuery.data as Entity).history[idx];

    data.push({
      diastolic: diastolic.state,
      systolic: systolic.state,
    });
  });

  return (
    <div className="p-4 grid grid-rows-2">
      <div className="row-start-1 row-span-1 flex flex-row justify-between">
        <div>
          <h1 className="text-2xl text-white">
            Trykket ser{" "}
            <span className="text-secondary-text font-bold">fint</span> ut
          </h1>

          <h1>
            {data[0].systolic}/{data[0].diastolic} mmHg
          </h1>
        </div>
        <div className="animate-pulse">
          <AiTwotoneHeart size={48} color={Scheme.Red} />
        </div>
      </div>
      <div>
        <ResponsiveContainer height={160}>
          <AreaChart
            data={data}
            margin={{
              top: 10,
              right: 30,
              left: 0,
              bottom: 0,
            }}
          >
            <Area
              type="monotone"
              dataKey="systolic"
              stroke="#8884d8"
              fill={Scheme.Red}
            />
            <Area
              type="monotone"
              dataKey="diastolic"
              stroke="#8884d8"
              fill={Scheme.White}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};
