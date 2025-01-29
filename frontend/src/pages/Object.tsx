import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import { ArrowLeft, Plus, Package, Filter } from 'lucide-react';
import { Card, CardHeader, CardDescription } from "@/components/ui/card";
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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { objectService } from '../services/object';
import { Object } from '../types/object';
import { authService } from '../services/auth';

const objectTypes = [
  'Furniture',
  'Electronics',
  'Appliances',
  'Decorations',
  'Tools',
  'Other'
];

export default function ObjectsPage() {
  const { id } = useParams<{ id?: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const roomName = location.state?.roomName || 'Room';
  const currentUser = authService.getCurrentUser();
  const [usernames, setUsernames] = useState<Record<string, string>>({});

  const [objects, setObjects] = useState<Object[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [newObjectName, setNewObjectName] = useState('');
  const [newObjectType, setNewObjectType] = useState(objectTypes[0]);
  const [isCreating, setIsCreating] = useState(false);
  const [showOnlyAvailable, setShowOnlyAvailable] = useState(false);

  const fetchUserInfo = async (userId: string) => {
    try {
      const user = await authService.getUser(userId);
      setUsernames(prev => ({
        ...prev,
        [userId]: user.username
      }));
    } catch (err) {
      console.error('Failed to fetch user info:', err);
    }
  };


  const fetchObjects = async () => {
    try {
      let response;
      if (id) {
        response = await objectService.listObjects(id);
      } else {
        response = await objectService.listObjects();
      }
      setObjects(response.data);
      
      const reservedObjects = response.data.filter(obj => obj.isReserved);
      reservedObjects.forEach(obj => {
        if (obj.reservedBy && !usernames[obj.reservedBy]) {
          fetchUserInfo(obj.reservedBy);
        }
      });
      
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchObjects();
  }, []);

  const handleCreateObject = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) {
      setError('Cannot create object: Room ID is missing');
      return;
    }

    setIsCreating(true);
    setError(null);

    try {
      await objectService.createObject({
        name: newObjectName,
        type: newObjectType,
        room_id: id
      });
      setNewObjectName('');
      setNewObjectType(objectTypes[0]);
      setIsDialogOpen(false);
      fetchObjects();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsCreating(false);
    }
  };

  const handleReserve = async (objectId: string) => {
    if (!id) {
      setError('Cannot reserve object: Room ID is missing');
      return;
    }

    if (!currentUser?.id) {
      navigate('/login');
      return;
    }

    try {
      await objectService.reserveObject(objectId, currentUser.id, id);
      fetchObjects();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const filteredObjects = showOnlyAvailable
    ? objects.filter(obj => !obj.isReserved)
    : objects;

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
            <Package className="h-8 w-8" />
            {roomName}
          </h1>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => setShowOnlyAvailable(!showOnlyAvailable)}
            className="flex items-center gap-2"
          >
            <Filter className="h-4 w-4" />
            {showOnlyAvailable ? 'Show All' : 'Show Available'}
          </Button>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                Add Object
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Add New Object to {roomName}</DialogTitle>
              </DialogHeader>
              <form onSubmit={handleCreateObject} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Object Name</Label>
                  <Input
                    id="name"
                    placeholder="Enter object name"
                    value={newObjectName}
                    onChange={(e) => setNewObjectName(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="type">Object Type</Label>
                  <Select value={newObjectType} onValueChange={setNewObjectType}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        {objectTypes.map((type) => (
                          <SelectItem key={type} value={type}>
                            {type}
                          </SelectItem>
                        ))}
                      </SelectGroup>
                    </SelectContent>
                  </Select>
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
                  {isCreating ? 'Adding...' : 'Add Object'}
                </Button>
              </form>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 text-red-600 p-4 rounded-md">
          {error}
        </div>
      )}

      {isLoading ? (
        <div className="text-center py-4">Loading objects...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredObjects.map((object) => (
            <Card
              key={object.id}
              className={object.isReserved ? 'bg-gray-50 border-yellow-200' : ''}
            >
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">{object.name}</h3>
                    <CardDescription>{object.type}</CardDescription>
                  </div>
                  <Package className={`h-5 w-5 ${object.isReserved ? 'text-yellow-500' : 'text-gray-500'}`} />
                </div>
                {object.isReserved ? (
                  <div className="space-y-2">
                    <div className="text-sm text-yellow-600 bg-yellow-50 p-2 rounded-md">
                      Reserved by: {(object.reservedBy && usernames[object.reservedBy]) || 'Loading...'}
                    </div>
                    {currentUser?.username === (object.reservedBy && usernames[object.reservedBy]) && (
                      <Button
                        variant="outline"
                        className="w-full border-red-200 text-red-600 hover:bg-red-50 hover:text-red-700"
                        onClick={async () => {
                          try {
                            await objectService.unreserveObject(object.id);
                            fetchObjects();
                          } catch (err: any) {
                            setError(err.message);
                          }
                        }}
                      >
                        Cancel Reservation
                      </Button>
                    )}
                  </div>
                ) : (
                  <Button
                    variant="outline"
                    className="w-full hover:bg-yellow-50 hover:text-yellow-600 hover:border-yellow-200"
                    onClick={() => handleReserve(object.id)}
                  >
                    Reserve
                  </Button>
                )}
              </CardHeader>
            </Card>
          ))}
        </div>
      )}

      {!isLoading && filteredObjects.length === 0 && (
        <div className="text-center p-8 border-2 border-dashed rounded-lg">
          <p className="text-gray-500">
            {showOnlyAvailable
              ? 'No available objects found. All objects are currently reserved.'
              : 'No objects found. Add your first object!'}
          </p>
        </div>
      )}
    </div>
  );
}