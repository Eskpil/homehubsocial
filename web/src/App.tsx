import {
  QueryClient,
  QueryClientProvider,
  useQuery,
} from "@tanstack/react-query";
import { Patients } from "./pages/Patients";
import { Patient } from "./pages/Patient";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Poll } from "./Poll";

const queryClient = new QueryClient();

const router = createBrowserRouter([
  {
    path: "/",
    element: <Patients />,
  },
  {
    path: "/patient/:id",
    element: <Patient />,
  },
]);

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="bg-primary-background overflow-hidden">
        <RouterProvider router={router} />
        <ToastContainer />
      </div>
    </QueryClientProvider>
  );
}

export default App;
