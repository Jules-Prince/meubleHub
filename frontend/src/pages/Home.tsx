import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Home } from '../types/home';
import { homeService } from '../services/home';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Plus, Home as HomeIcon } from 'lucide-react';

export default function HomePage() {
  const navigate = useNavigate(); // Add this line to get the navigate function
  const [homes, setHomes] = useState<Home[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newHomeName, setNewHomeName] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const fetchHomes = async () => {
    try {
      const response = await homeService.listHomes();
      setHomes(response.data);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchHomes();
  }, []);

  const handleCreateHome = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsCreating(true);
    setError(null);

    try {
      await homeService.createHome({ name: newHomeName });
      setNewHomeName('');
      setIsDialogOpen(false);
      fetchHomes(); // Refresh the list
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsCreating(false);
    }
  };

  if (isLoading) {
    return <div className="flex justify-center items-center p-8">Loading homes...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold">Homes</h1>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Add Home
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create New Home</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateHome} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Home Name</Label>
                <Input
                  id="name"
                  placeholder="Enter home name"
                  value={newHomeName}
                  onChange={(e) => setNewHomeName(e.target.value)}
                  required
                />
              </div>
              {error && (
                <div className="text-red-500 text-sm">
                  {error}
                </div>
              )}
              <Button 
                type="submit" 
                className="w-full"
                disabled={isCreating}
              >
                {isCreating ? 'Creating...' : 'Create Home'}
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 p-4 rounded-md">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {homes.map((home) => (
          <Card key={home.id}>
            <CardHeader className="flex flex-row items-center gap-2">
              <HomeIcon className="h-6 w-6" />
              <div>
                <CardTitle>{home.name}</CardTitle>
                <CardDescription>Home ID: {home.id}</CardDescription>
              </div>
            </CardHeader>
            <CardContent>
              <Button 
                variant="outline" 
                className="w-full" 
                onClick={() => navigate(`/homes/${home.id}/rooms`, {
                  state: { homeName: home.name }
                })}
              >
                View Rooms
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>

      {homes.length === 0 && !error && (
        <div className="text-center p-8 border-2 border-dashed rounded-lg">
          <p className="text-gray-500">No homes found. Create your first home!</p>
        </div>
      )}
    </div>
  );
}