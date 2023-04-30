import { useNavigate, useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { Patient } from "../types/Patient";
import axios from "axios";
import { HeartbeatChart } from "../components/HeartbeatChart";
import { useElementSize } from "usehooks-ts";
import { BloodpressureGraph } from "../components/BloodpressureChart";
import { toast } from "react-toastify";

interface Props {}

export const Patient: React.FC<Props> = () => {
  const { id } = useParams();
  const { data, error, isFetching } = useQuery<Patient>({
    queryKey: [`patient-${id}`],
    queryFn: async () =>
      axios
        .get(`http://localhost:8080/patients/${id}/`)
        .then((res) => res.data),
  });

  const navigate = useNavigate();

  if (isFetching) {
    return <h1>loading...</h1>;
  }

  const onBackClick = () => {
    navigate("/");
  };

  const onRequestExaminationClick = () => {
    toast(`${data!.name} Er nå kalt inn til en undersøkelse`);
  };
  const onRequestRoutineCheck = () => {
    toast(`${data!.name} Er nå kalt inn til en rutinemessing sjekk`);
  };

  return (
    <div className="p-16">
      <div className="bg-secondary-background rounded rounded-lg p-4">
        <button
          className="bg-accent flex flex-col justify-center p-4 rounded rounded-lg pr-12 pl-12 hover:cursor-pointer hover:opacity-90"
          onClick={() => onBackClick()}
        >
          <h1 className="text-accent-text font-bold">Tilbake</h1>
        </button>
      </div>
      <div className="grid grid-cols-4 gap-8 grid-rows-3 mt-16">
        <div className="col-span-1 col-start-1 h-40">
          <div className="bg-secondary-background rounded rounded-lg">
            <img
              src="/edith-full.jpeg"
              className="rounded rounded-lg"
              alt={`Bilde av pasienten ${data!.name}`}
            />
            <div className="p-4">
              <div className="text-2xl">
                <h1 className="text-white">
                  Edith,{" "}
                  <span className="font-bold text-secondary-text">
                    {data!.age} år
                  </span>
                </h1>
              </div>
              <div className="mt-4">
                <p className="text-white">
                  Her kan du <span className="font-bold">monitorere</span> og{" "}
                  <span className="font-bold">administere</span> din pasient{" "}
                  <span className="font-bold">{data!.name}</span>
                </p>
              </div>
            </div>
          </div>
        </div>
        <div className="col-span-2 col-start-2 rounded rounded-lg row-span-1 row-start-1">
          <div className="bg-secondary-background rounded rounded-lg h-max">
            <HeartbeatChart
              patient={data! as Patient}
              width={1200}
              height={160}
            />
          </div>
        </div>
        <div className="col-span-2 col-start-1 rounded rounded-lg row-span-1 row-start-2">
          <div className="bg-secondary-background rounded rounded-lg h-max">
            <BloodpressureGraph patient={data! as Patient} />
          </div>
        </div>
        <div className="row-start-2 col-start-3 h-40">
          <div className="bg-secondary-background rounded rounded-lg">
            <div className="p-4 flex flex-col justify-around">
              <button
                className="bg-accent flex flex-col justify-center p-4 rounded rounded-lg pr-12 pl-12 hover:cursor-pointer hover:opacity-90"
                onClick={() => onRequestExaminationClick()}
              >
                <h1 className="text-accent-text font-bold">
                  Kall inn {data!.name} til en undersøkelse
                </h1>
              </button>
              <button
                className="mt-4 bg-accent flex flex-col justify-center p-4 rounded rounded-lg pr-12 pl-12 hover:cursor-pointer hover:opacity-90"
                onClick={() => onRequestRoutineCheck()}
              >
                <h1 className="text-accent-text font-bold">
                  Kall inn {data!.name} til en rutinemessing sjekk
                </h1>
              </button>
            </div>
          </div>
        </div>

        <div className="col-span-1 col-start-4 row-span-2 row-start-1">
          <div className="bg-secondary-background h-max rounded rounded-lg">
            <div className="p-8">
              <img src="/logo.png" alt="" className="rounded rounded-lg" />
            </div>
            <div className="p-8">
              <h1 className="text-xl text-white">
                <span className="font-bold text-secondary-text">
                  Home Hub Social
                </span>
                , Et produkt laget av Linus, Tage, Erik og Noah.
              </h1>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
