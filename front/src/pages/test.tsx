import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export default function TestPage() {
  return (
    <Alert variant="default">
      <AlertTitle>Heads up!</AlertTitle>
      <AlertDescription>
        You can add components and dependencies to your app using the cli.
      </AlertDescription>
    </Alert>
  );
}
