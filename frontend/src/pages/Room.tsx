import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import { Home as HomeIcon, ArrowLeft, Plus, DoorClosed, Trash2 } from 'lucide-react';
import { Card, CardHeader, CardContent, CardDescription } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { roomService } from '../services/room';
import { Room } from '../types/room';
import { authService } from '@/services/auth';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';

export default function HomeRooms() {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const homeId = Number(id);
  const homeName = location.state?.homeName || 'Home';

  const [rooms, setRooms] = useState<Room[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newRoomName, setNewRoomName] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [roomToDelete, setRoomToDelete] = useState<number | null>(null);
  const currentUser = authService.getCurrentUser();

  const fetchRooms = async () => {
    try {
      const response = await roomService.listRooms(homeId);
      setRooms(response.data);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRoom = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsCreating(true);
    setError(null);

    try {
      await roomService.createRoom({
        name: newRoomName,
        home_id: homeId
      });
      setNewRoomName('');
      setIsDialogOpen(false);
      fetchRooms();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsCreating(false);
    }
  };

  useEffect(() => {
    if (homeId) {
      fetchRooms();
    }
  }, [homeId]);

  const handleDeleteRoom = async () => {
    if (!roomToDelete) return;

    try {
      await roomService.deleteRoom(roomToDelete);
      setRoomToDelete(null);
      fetchRooms();
    } catch (err: any) {
      setError(err.message);
    }
  };


  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate(-1)}
          >
            <ArrowLeft className="h-6 w-6" />
          </Button>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <HomeIcon className="h-8 w-8" />
            {homeName} - Rooms
          </h1>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              Add Room
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Add Room to {homeName}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateRoom} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="name">Room Name</Label>
                <Input
                  id="name"
                  placeholder="Enter room name"
                  value={newRoomName}
                  onChange={(e) => setNewRoomName(e.target.value)}
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
                {isCreating ? 'Creating...' : 'Create Room'}
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

      {isLoading ? (
        <div className="text-center py-4">Loading rooms...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {rooms.map((room) => (
            <Card key={room.id}>
              <CardHeader className="flex flex-row items-center gap-2">
                <DoorClosed className="h-5 w-5" />
                <div>
                  <h3 className="font-semibold">{room.name}</h3>
                  <CardDescription>Room ID: {room.id}</CardDescription>
                </div>
                {currentUser?.isAdmin && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-red-500 hover:text-red-700 hover:bg-red-50"
                    onClick={() => setRoomToDelete(room.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </CardHeader>
              <CardContent>
                <Button
                  variant="outline"
                  className="w-full"
                  onClick={() => navigate(`/room/${room.id}`, {
                    state: { roomName: room.name }
                  })}
                >
                  View Details
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      <AlertDialog open={!!roomToDelete} onOpenChange={(open) => !open && setRoomToDelete(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the room
              and all its associated rooms and objects.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteRoom}
              className="bg-red-500 hover:bg-red-600"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {!isLoading && rooms.length === 0 && (
        <div className="text-center p-8 border-2 border-dashed rounded-lg">
          <p className="text-gray-500">No rooms found. Create your first room!</p>
        </div>
      )}
    </div>
  );
}