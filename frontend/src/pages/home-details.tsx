import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ArrowLeft, Home as HomeIcon } from 'lucide-react';
import { Card, CardHeader, CardContent } from "@/components/ui/card";
import { Button } from '@/components/ui/button';

export default function HomeDetailsPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  return (
    <div className="space-y-6">
        <Button 
          variant="ghost" 
          size="icon"
          onClick={() => navigate('/home')}
        >
          <ArrowLeft className="h-6 w-6" />
        </Button>
      <h1 className="text-3xl font-bold flex items-center gap-2">
        <HomeIcon className="h-8 w-8" />
        Home Details
      </h1>
      
      <Card>
        <CardHeader>
          <h2 className="text-xl font-semibold">Home ID: {id}</h2>
        </CardHeader>
        <CardContent>
          {/* Add more home details here */}
          <p className="text-gray-500">More details coming soon...</p>
        </CardContent>
      </Card>
    </div>
  );
}