package scalable.p2.thumbnail.api;

import org.springframework.http.ResponseEntity;
import scalable.p2.thumbnail.repository.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import scalable.p2.thumbnail.statuses.Statuses;

import java.io.IOException;
import java.util.List;
import java.util.Map;

@RestController
public class APIController {
    private final APIService service;

    @Autowired
    public APIController(APIService service) {
        this.service = service;
    }

    @RequestMapping(value = "/api/submit", method = RequestMethod.POST)
    public void submitJob(@RequestBody Pojo payload) {
        List<Video> list = payload.getVideos();

        for (Video v : list) {
            try {
                this.service.queue(payload.getBucketID(), v.getVideoName(), v.getDuration(), v.getFPS());
            } catch (Exception e) {
                ResponseEntity.badRequest();
            }
        }
    }

   @RequestMapping(value = "api/status/extract", method = RequestMethod.GET)
    public Map<String, Integer> getExtractStatus(@RequestParam Integer BucketID, @RequestParam String VideoName) {
        return this.service.getStatus(BucketID, VideoName, Statuses.EXTRACT);
   }

    @RequestMapping(value = "api/status/resize", method = RequestMethod.GET)
    public Map<String, Integer> getResizeStatus(@RequestParam Integer BucketID, @RequestParam String VideoName) {
        return this.service.getStatus(BucketID, VideoName, Statuses.RESIZE);
    }

    @RequestMapping(value = "api/status/compile", method = RequestMethod.GET)
    public Map<String, Integer> getCompileStatus(@RequestParam Integer BucketID, @RequestParam String VideoName) {
        return this.service.getStatus(BucketID, VideoName, Statuses.COMPILE);
    }

    @RequestMapping(value = "api/get/gifs", method = RequestMethod.GET)
    public List<Status> getGifs(@RequestParam Integer BucketID) {
        return this.service.getGifs(BucketID);
    }
}