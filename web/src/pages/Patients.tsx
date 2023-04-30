import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { Patient } from "../types/Patient";
import { useNavigate } from "react-router-dom";

interface Props {}

export const Patients: React.FC<Props> = () => {
  const { data, error, isFetching } = useQuery<Patient[]>({
    queryKey: ["patients"],
    queryFn: async () =>
      axios.get("http://localhost:8080/patients/").then((res) => res.data),
  });

  const navigate = useNavigate();

  if (isFetching) {
    return <div>loading...</div>;
  }

  const onProfileClick = (patient: Patient) => {
    return navigate(`/patient/${patient.id}`);
  };

  return (
    <div className="p-20">
      <div className="grid grid-cols-6">
        <div className="col-span-1 col-start-1">
          <div className="bg-secondary-background rounded rounded-lg">
            <img
              src="/nurse.jpg"
              className="rounded rounded-lg"
              alt="Bilde av sykepleieren Daniel"
            />
            <div className="p-4">
              <div className="text-2xl">
                <h1 className="text-white">
                  Hei,{" "}
                  <span className="font-bold text-secondary-text">Daniel</span>
                </h1>
              </div>
              <div className="mt-4">
                <p className="text-white">
                  Her kan du <span className="font-bold">monitorere</span> og{" "}
                  <span className="font-bold">administere</span> dine pasienter
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="col-span-3 col-start-3">
          {data!.map((patient) => (
            <div
              key={patient.id}
              className="rounded rounded-lg p-4 bg-secondary-background flex flex-row justify-between"
            >
              <div className="flex flex-row">
                <img
                  src="/edith-cropped.png"
                  width={80}
                  height={80}
                  className="rounded rounded-lg"
                  alt="Edith's bilde"
                />
                <div className="ml-12 flex flex-col justify-center">
                  <h1 className="text-white">
                    {patient.name},{" "}
                    <span className="text-secondary-text">
                      {patient.age} år
                    </span>
                  </h1>
                </div>
              </div>

              <button
                className="bg-accent flex flex-col justify-center p-4 rounded rounded-lg pr-12 pl-12 hover:cursor-pointer hover:opacity-90"
                onClick={() => onProfileClick(patient)}
              >
                <h1 className="text-accent-text font-bold">Se profil</h1>
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
