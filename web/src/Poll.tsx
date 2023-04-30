import { useQuery } from "@tanstack/react-query";
import { Event } from "./types/Event";
import axios from "axios";
import { toast } from "react-toastify";

interface Props {}

export const Poll: React.FC<Props> = () => {
  const key = "events";

  const { data, isLoading } = useQuery<Event[]>({
    queryKey: [key],
    queryFn: axios.get("http://localhost:8080/events/").then((res) => res.data),
    refetchInterval: () => 500,
  });

  if (!isLoading) {
    (data! as Event[]).forEach((event) => {
      toast(event.content);
    });
  }

  return <></>;
};
